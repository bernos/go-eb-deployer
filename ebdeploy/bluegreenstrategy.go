package ebdeploy

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/s3"
	//	"io"
	//"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	//"time"
)

func NewBlueGreenStrategy() *DeploymentPipeline {
	pipeline := new(DeploymentPipeline)
	pipeline.AddStep(ensureBucketExists)
	pipeline.AddStep(uploadVersion)
	pipeline.AddStep(prepareTargetEnvironment)
	pipeline.AddStep(deployApplicationVersion)
	pipeline.AddStep(runSmokeTest)

	return pipeline
}

func ensureBucketExists(ctx *DeploymentContext, next Continue) error {

	var (
		bucket string
		exists bool
		err    error
	)

	s3Service := services.NewS3Service(s3.New(ctx.AwsConfig))

	if bucket, err = ctx.Bucket(); err != nil {
		return err
	}

	if exists, err = s3Service.BucketExists(bucket); err != nil {
		return err
	}

	if exists {
		log.Printf("Bucket %s already exists", bucket)
		return next()
	}

	log.Printf("Creating bucket %s", bucket)

	if err = s3Service.CreateBucket(bucket); err != nil {
		return err
	}

	log.Printf("Created bucket %s", bucket)

	return next()
}

func uploadVersion(ctx *DeploymentContext, next Continue) error {

	var (
		bucket        string
		err           error
		versionExists bool
	)

	s3Service := services.NewS3Service(s3.New(ctx.AwsConfig))
	ebService := services.NewEBService(elasticbeanstalk.New(ctx.AwsConfig))
	key := ctx.Version + ".zip"

	if versionExists, err = ebService.ApplicationVersionExists(ctx.Configuration.ApplicationName, ctx.Version); err != nil {
		return err
	} else if versionExists {
		return errors.New("Version " + ctx.Version + " already exists")
	}

	if bucket, err = ctx.Bucket(); err != nil {
		return err
	}

	if err = s3Service.UploadFile(ctx.SourceBundle, bucket, key); err != nil {
		return err
	}

	log.Printf("Uploaded version %s", ctx.Version)

	if err = ebService.CreateApplicationVersion(ctx.Configuration.ApplicationName, ctx.Version, bucket, key); err != nil {
		return err
	}

	log.Printf("Created version %s", ctx.Version)

	return next()
}

func prepareTargetEnvironment(ctx *DeploymentContext, next Continue) error {

	ebService := services.NewEBService(elasticbeanstalk.New(ctx.AwsConfig))

	if environments, err := ebService.GetEnvironments(ctx.Configuration.ApplicationName); err == nil {

		log.Printf("Found %d existing environments", len(environments))
		requiresTerminate := false
		activeCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, true)
		inactiveCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, false)

		activeEnvironment := findEnvironment(environments, cnamePredicate(activeCname))
		inactiveEnvironment := findEnvironment(environments, cnamePredicate(inactiveCname))

		if activeEnvironment != nil && inactiveEnvironment != nil {
			log.Println("Both active and inactive environments were found. Inactive environment will be terminated.")

			ctx.TargetEnvironment = &TargetEnvironment{
				Name:     *inactiveEnvironment.EnvironmentName,
				CNAME:    inactiveCname,
				IsActive: false,
			}

			requiresTerminate = true
		} else if activeEnvironment == nil && inactiveEnvironment == nil {
			log.Println("Neither active nor inactive environments were found. Deploying directly to active environment")

			ctx.TargetEnvironment = &TargetEnvironment{
				Name:     calculateEnvironmentName(ctx.Environment, "a"),
				CNAME:    activeCname,
				IsActive: true,
			}
		} else if activeEnvironment != nil {
			activeSuffix := getSuffixFromEnvironmentName(*activeEnvironment.EnvironmentName)
			inactiveSuffix := "a"

			if activeSuffix == "a" {
				inactiveSuffix = "b"
			}

			log.Printf("Active environment '%s' found. Deploying to inactive environment '%s'", activeSuffix, inactiveSuffix)

			ctx.TargetEnvironment = &TargetEnvironment{
				Name:     calculateEnvironmentName(ctx.Environment, inactiveSuffix),
				CNAME:    inactiveCname,
				IsActive: false,
			}
		} else {
			return errors.New("Current environment state is not recognized. Please there is either one recognised, active environment, or no environments at all")
		}

		ctx.TargetEnvironment.Url = "http://" + ctx.TargetEnvironment.CNAME + ".elasticbeanstalk.com"

		done := make(chan struct{})
		defer close(done)
		ebService.LogEnvironmentEvents(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, done)

		if requiresTerminate {
			if err := ebService.TerminateEnvironment(*inactiveEnvironment.EnvironmentID); err != nil {
				return err
			}

			if err := ebService.WaitForEnvironment(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, func(e *elasticbeanstalk.EnvironmentDescription) bool {
				return *e.Status == "Terminated"
			}); err != nil {
				return err
			}
		}

		log.Printf("active environment: %v", activeEnvironment)
		log.Printf("inactive environment: %v", inactiveEnvironment)
		log.Printf("Target environment: %v", ctx.TargetEnvironment)

		return next()
	} else {
		return err
	}
}

func deployApplicationVersion(ctx *DeploymentContext, next Continue) error {

	log.Printf("Deploying version %s to environment %s", ctx.Version, ctx.TargetEnvironment.Name)

	client := elasticbeanstalk.New(ctx.AwsConfig)
	ebService := services.NewEBService(client)
	params := &elasticbeanstalk.CreateEnvironmentInput{
		ApplicationName:   aws.String(ctx.Configuration.ApplicationName),
		EnvironmentName:   aws.String(ctx.TargetEnvironment.Name),
		CNAMEPrefix:       aws.String(ctx.TargetEnvironment.CNAME),
		OptionSettings:    []*elasticbeanstalk.ConfigurationOptionSetting(ctx.Configuration.OptionSettings),
		Tags:              []*elasticbeanstalk.Tag(ctx.Configuration.Tags),
		SolutionStackName: aws.String(ctx.Configuration.SolutionStackName),
		Tier:              ctx.Configuration.Tier,
		VersionLabel:      aws.String(ctx.Version),
	}

	if resp, err := client.CreateEnvironment(params); err != nil {
		return err
	} else {

		log.Printf("Deployed: %v", resp)
		log.Printf("Waiting for environment health to go green")

		if err := ebService.WaitForEnvironment(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, func(e *elasticbeanstalk.EnvironmentDescription) bool {
			return *e.Health == "Green"
		}); err != nil {
			return err
		} else {
			log.Printf("Environment is green")
			return next()
		}
	}
}

func runSmokeTest(ctx *DeploymentContext, next Continue) error {
	if len(ctx.Configuration.SmokeTestUrl) > 0 {
		url := strings.Replace(ctx.Configuration.SmokeTestUrl, "{url}", ctx.TargetEnvironment.Url, -1)

		log.Printf("Running smoke test against %s", url)

		for {

			if resp, err := http.Get(url); err != nil {
				return err
			} else if resp.StatusCode == 200 {
				log.Printf("Smoke test passed!")
				break
			}

			time.Sleep(5 * time.Second)
		}

	} else {
		log.Println("No SmokeTestUrl specified. Skipping smoke tests")
	}

	return next()
}

func cnamePredicate(cname string) func(*elasticbeanstalk.EnvironmentDescription) bool {
	return func(e *elasticbeanstalk.EnvironmentDescription) bool {
		return *e.CNAME == cname+".elasticbeanstalk.com"
	}
}

func calculateCnamePrefix(applicationName string, environmentName string, isActive bool) string {
	whiteSpaceToHyphenRegexp = regexp.MustCompile(`\s`)
	cname := strings.ToLower(whiteSpaceToHyphenRegexp.ReplaceAllString(applicationName+"-"+environmentName, "-"))

	if isActive {
		return cname
	}

	return cname + "-inactive"
}

func calculateEnvironmentName(name string, suffix string) string {
	return strings.ToLower(name + "-" + suffix + "-" + randomString(8))
}

func randomString(length int) string {
	u := make([]byte, length/2)
	_, err := rand.Read(u)

	if err != nil {
		return ""
	}

	return hex.EncodeToString(u)
}

func findEnvironment(environments []*elasticbeanstalk.EnvironmentDescription, predicate func(*elasticbeanstalk.EnvironmentDescription) bool) *elasticbeanstalk.EnvironmentDescription {
	for _, e := range environments {
		if predicate(e) {
			return e
		}
	}
	return nil
}

func getSuffixFromEnvironmentName(name string) string {
	tokens := strings.Split(name, "-")

	return tokens[len(tokens)-2]
}

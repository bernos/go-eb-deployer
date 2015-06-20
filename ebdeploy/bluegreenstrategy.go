package ebdeploy

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/s3"
	//	"io"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
	"regexp"
	"strings"
)

func NewBlueGreenStrategy() *DeploymentPipeline {
	pipeline := new(DeploymentPipeline)
	pipeline.AddStep(ensureBucketExists)
	pipeline.AddStep(uploadVersion)
	pipeline.AddStep(prepareTargetEnvironment)

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

	// TODO: return error if the version already exists

	var (
		bucket string
		err    error
	)

	s3Service := services.NewS3Service(s3.New(ctx.AwsConfig))
	ebClient := elasticbeanstalk.New(ctx.AwsConfig)
	key := ctx.Version + ".zip"

	if bucket, err = ctx.Bucket(); err != nil {
		return err
	}

	if err = s3Service.UploadFile(ctx.SourceBundle, bucket, key); err != nil {
		return err
	}

	log.Printf("Uploaded version %s", ctx.Version)

	if err = createApplicationVersion(ebClient, ctx.Configuration.ApplicationName, ctx.Version, bucket, key); err != nil {
		return err
	}

	log.Printf("Created version %s", ctx.Version)

	return next()
}

func prepareTargetEnvironment(ctx *DeploymentContext, next Continue) error {

	svc := elasticbeanstalk.New(&aws.Config{Region: ctx.Configuration.Region})

	if environments, err := getEnvironments(svc, ctx.Configuration.ApplicationName); err == nil {

		log.Printf("Found %d existing environments", len(environments))

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

			if err := terminateEnvironment(svc, *inactiveEnvironment.EnvironmentID); err != nil {
				return err
			}
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

			log.Printf("Active environment '%s' found. Deploying to inactive environment '%s'", activeSuffix, inactiveSuffix)

			if activeSuffix == "a" {
				inactiveSuffix = "b"
			}

			ctx.TargetEnvironment = &TargetEnvironment{
				Name:     calculateEnvironmentName(ctx.Environment, inactiveSuffix),
				CNAME:    inactiveCname,
				IsActive: false,
			}
		} else {
			return errors.New("Current environment state is not recognized. Please there is either one recognised, active environment, or no environments at all")
		}

		ctx.TargetEnvironment.Url = "http://" + ctx.TargetEnvironment.CNAME + ".elasticbeanstalk.com"

		log.Printf("active environment: %v", activeEnvironment)
		log.Printf("inactive environment: %v", inactiveEnvironment)

		log.Printf("Target environment: %v", ctx.TargetEnvironment)

		return next()
	} else {
		return err
	}
}

func getEnvironments(svc *elasticbeanstalk.ElasticBeanstalk, applicationName string) ([]*elasticbeanstalk.EnvironmentDescription, error) {
	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		ApplicationName: aws.String(applicationName),
		IncludeDeleted:  aws.Boolean(false),
	}

	if resp, err := svc.DescribeEnvironments(params); err == nil {
		return resp.Environments, nil
	} else {
		return nil, err
	}
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

func terminateEnvironment(svc *elasticbeanstalk.ElasticBeanstalk, environmentId string) error {
	params := &elasticbeanstalk.TerminateEnvironmentInput{
		EnvironmentID: &environmentId,
	}

	if _, err := svc.TerminateEnvironment(params); err != nil {
		return err
	}
	// TODO: Wait for env to be terminated...
	return nil
}

func createApplicationVersion(svc *elasticbeanstalk.ElasticBeanstalk, applicationName string, version string, bucket string, key string) error {
	params := &elasticbeanstalk.CreateApplicationVersionInput{
		ApplicationName:       aws.String(applicationName), // Required
		VersionLabel:          aws.String(version),         // Required
		AutoCreateApplication: aws.Boolean(true),
		SourceBundle: &elasticbeanstalk.S3Location{
			S3Bucket: aws.String(bucket),
			S3Key:    aws.String(key),
		},
	}

	if _, err := svc.CreateApplicationVersion(params); err != nil {
		return err
	}

	return nil
}

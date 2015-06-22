package bluegreen

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"regexp"
	"strings"
)

func init() {
	pipeline.RegisterStrategy("blue-green", NewBlueGreenStrategy)
}

func NewBlueGreenStrategy() *pipeline.DeploymentPipeline {
	pipeline := new(pipeline.DeploymentPipeline)
	pipeline.AddStep(ensureBucketExists)
	pipeline.AddStep(uploadVersion)
	pipeline.AddStep(prepareTargetEnvironment)
	pipeline.AddStep(deployApplicationVersion)
	pipeline.AddStep(runSmokeTest)
	pipeline.AddStep(swapCnames)
	return pipeline
}

func cnamePredicate(cname string) func(*elasticbeanstalk.EnvironmentDescription) bool {
	return func(e *elasticbeanstalk.EnvironmentDescription) bool {
		return *e.CNAME == cname+".elasticbeanstalk.com"
	}
}

func calculateCnamePrefix(applicationName string, environmentName string, isActive bool) string {
	whiteSpaceToHyphenRegexp := regexp.MustCompile(`\s`)
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

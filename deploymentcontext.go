package ebdeploy

import (
	"errors"
	"regexp"
	"strings"
)

var (
	whiteSpaceToHyphenRegexp = regexp.MustCompile(`\s`)
)

type DeploymentContext struct {
	Configuration *Configuration
	Environment   string
	SourceBundle  string
	Version       string
}

func (d *DeploymentContext) Bucket() (string, error) {
	if d.Configuration == nil {
		return "", errors.New("DeploymentContext has no Configuration")
	}

	if len(d.Configuration.Bucket) == 0 && len(d.Configuration.ApplicationName) == 0 {
		return "", errors.New("Either Configuration.Bucket or Configuration.ApplicationName are required")
	}

	if len(d.Configuration.Bucket) > 0 {
		return d.Configuration.Bucket, nil
	}

	return calculateBucketName(d.Configuration.ApplicationName), nil
}

func calculateBucketName(s string) string {
	return strings.ToLower(whiteSpaceToHyphenRegexp.ReplaceAllString(s, "-") + "-packages")
}

func NewDeploymentContext(configuration *Configuration, environment string, sourceBundle string, version string) (*DeploymentContext, error) {

	if configuration == nil {
		return nil, errors.New("Configuration is required")
	}

	if !configuration.HasEnvironment(environment) {
		return nil, errors.New("Invalid environment " + environment)
	}

	if len(version) == 0 {
		return nil, errors.New("Invalid version number")
	}

	if len(sourceBundle) == 0 {
		return nil, errors.New("Invalid version number")
	}

	d := &DeploymentContext{
		Configuration: configuration,
		Environment:   environment,
		SourceBundle:  sourceBundle,
		Version:       version,
	}

	return d, nil
}

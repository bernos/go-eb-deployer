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

func NewDeploymentContext(c *Configuration) *DeploymentContext {
	d := new(DeploymentContext)
	d.Configuration = c
	return d
}

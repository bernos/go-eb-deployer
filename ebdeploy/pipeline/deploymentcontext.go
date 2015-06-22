package pipeline

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/bernos/go-eb-deployer/ebdeploy/config"
)

var (
	whiteSpaceToHyphenRegexp = regexp.MustCompile(`\s`)
)

type TargetEnvironment struct {
	Name     string
	CNAME    string
	IsActive bool
	Url      string
}

type DeploymentContext struct {
	Configuration     *config.Configuration
	Environment       string
	SourceBundle      string
	Version           string
	TargetEnvironment *TargetEnvironment
	AwsConfig         *aws.Config
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

func NewDeploymentContext(configuration *config.Configuration, environment string, sourceBundle string, version string) (*DeploymentContext, error) {

	if configuration == nil {
		return nil, errors.New("Configuration is required")
	}

	if !configuration.HasEnvironment(environment) {
		return nil, errors.New("Invalid environment " + environment)
	}

	if len(sourceBundle) == 0 {
		return nil, errors.New("Invalid source bundle")
	}

	if len(version) == 0 {
		v, err := calculateVersionFromFileContents(sourceBundle)

		if err != nil {
			return nil, err
		}

		version = v
	}

	d := &DeploymentContext{
		Configuration: configuration,
		Environment:   environment,
		SourceBundle:  sourceBundle,
		Version:       version,
		AwsConfig: &aws.Config{
			Region: configuration.Region,
		},
	}

	return d, nil
}

func calculateVersionFromFileContents(file string) (string, error) {
	hasher := md5.New()

	if f, err := os.Open(file); err == nil {
		defer f.Close()

		if _, err := io.Copy(hasher, f); err == nil {
			return hex.EncodeToString(hasher.Sum(nil)), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

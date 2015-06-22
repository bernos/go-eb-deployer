package pipeline

import (
	"github.com/bernos/go-eb-deployer/ebdeploy/config"
	"testing"
)

func TestCalculateBucketName(t *testing.T) {
	expected := "my-application-name-packages"
	actual := calculateBucketName("My Application name")

	if expected != actual {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestNewDeploymentContext(t *testing.T) {
	config := new(config.Configuration)

	if _, err := NewDeploymentContext(config, "staging", "test.zip", "1.0.1"); err == nil {
		t.Error("Expected invalid environment error")
	}
}

func TestBucket(t *testing.T) {
	cfg := new(config.Configuration)
	cfg.Environments = []config.Environment{config.Environment{Name: "staging"}}
	context, _ := NewDeploymentContext(cfg, "staging", "test.zip", "1.0.1")

	expected := "my-application-name-packages"

	if _, err := context.Bucket(); err == nil {
		t.Errorf("Expected err but got nil")
	}

	cfg.ApplicationName = "My Application Name"

	if actual, err := context.Bucket(); err == nil {
		if actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	} else {
		t.Errorf("%v", err)
	}

	cfg.Bucket = "my-bucket"

	if actual, err := context.Bucket(); err == nil {
		if actual != cfg.Bucket {
			t.Errorf("Expected %s, but got %s", cfg.Bucket, actual)
		}
	} else {
		t.Errorf("%v", err)
	}
}

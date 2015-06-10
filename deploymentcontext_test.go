package ebdeploy

import (
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
	config := new(Configuration)

	if _, err := NewDeploymentContext(config, "staging", "test.zip", "1.0.1"); err == nil {
		t.Error("Expected invalid environment error")
	}
}

func TestBucket(t *testing.T) {
	config := new(Configuration)
	config.Environments = []Environment{Environment{Name: "staging"}}
	context, _ := NewDeploymentContext(config, "staging", "test.zip", "1.0.1")

	expected := "my-application-name-packages"

	if _, err := context.Bucket(); err == nil {
		t.Errorf("Expected err but got nil")
	}

	config.ApplicationName = "My Application Name"

	if actual, err := context.Bucket(); err == nil {
		if actual != expected {
			t.Errorf("Expected %s, but got %s", expected, actual)
		}
	} else {
		t.Errorf("%v", err)
	}

	config.Bucket = "my-bucket"

	if actual, err := context.Bucket(); err == nil {
		if actual != config.Bucket {
			t.Errorf("Expected %s, but got %s", config.Bucket, actual)
		}
	} else {
		t.Errorf("%v", err)
	}
}

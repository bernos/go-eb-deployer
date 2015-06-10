package ebdeploy

import (
	"testing"
)

func TestBucketExists(t *testing.T) {
	exists, err := bucketExists("my-application-packages", "ap-southeast-2")

	if err != nil {
		t.Errorf("Test failed %s", err)
	}

	if !exists {
		t.Errorf("Expected bucket to exist")
	}
}

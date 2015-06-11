package ebdeploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"testing"
)

func IgnoreTestBucketExists(t *testing.T) {

	svc := s3.New(&aws.Config{Region: "ap-southeast-2"})

	exists, err := bucketExists(svc, "my-application-packages")

	if err != nil {
		t.Errorf("Test failed %s", err)
	}

	if !exists {
		t.Errorf("Expected bucket to exist")
	}
}

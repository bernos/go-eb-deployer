package ebdeploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

func Deploy(ctx *DeploymentContext) {
	//	ensureBucketExists(ctx.Bucket())
}

func bucketExists(bucket string, region string) (bool, error) {
	svc := s3.New(&aws.Config{Region: region})

	if output, err := svc.ListBuckets(new(s3.ListBucketsInput)); err == nil {
		for _, b := range output.Buckets {
			if *b.Name == bucket {
				return true, nil
			}
		}
		return false, nil
	} else {
		return false, err
	}
}

func ensureBucketExists(bucket string, region string) error {
	svc := s3.New(&aws.Config{Region: region})

	if output, err := svc.ListBuckets(new(s3.ListBucketsInput)); err == nil {
		for _, bucket := range output.Buckets {
			log.Printf("bucket: %s", *bucket.Name)
		}
		return nil
	} else {
		return err
	}
}

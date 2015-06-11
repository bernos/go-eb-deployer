package ebdeploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	//	"log"
)

type DeploymentStep func(*DeploymentContext, *DeploymentStep)

func Deploy(ctx *DeploymentContext) {
	//	ensureBucketExists(ctx.Bucket())

	svc := s3.New(&aws.Config{Region: ctx.Configuration.Region})

	if bucket, err := ctx.Bucket(); err == nil {
		if ensureBucketExists(svc, bucket) == nil {

		}
	}
}

func bucketExists(svc *s3.S3, bucket string) (bool, error) {
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

func createBucket(svc *s3.S3, bucket string) error {
	return nil
}

func ensureBucketExists(svc *s3.S3, bucket string) error {
	if exists, err := bucketExists(svc, bucket); err == nil {
		if !exists {
			return createBucket(svc, bucket)
		}
		return nil
	} else {
		return err
	}
}

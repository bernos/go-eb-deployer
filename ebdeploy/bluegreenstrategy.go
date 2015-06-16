package ebdeploy

import (
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

func NewBlueGreenStrategy() *DeploymentPipeline {
	pipeline := new(DeploymentPipeline)
	pipeline.AddStep(ensureBucketExists)
	pipeline.AddStep(uploadVersion)

	return pipeline
}

func ensureBucketExists(ctx *DeploymentContext, next Continue) error {
	log.Printf("Ensure bucket exists")
	return next()
}

func uploadVersion(ctx *DeploymentContext, next Continue) error {
	log.Printf("Upload version")
	return next()
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

/*
func ensureBucketExists(svc *s3.S3, bucket string) error {
	if exists, err := bucketExists(svc, bucket); err == nil {
		if !exists {
			return createBucket(svc, bucket)
		}
		return nil
	} else {
		return err
	}
}*/

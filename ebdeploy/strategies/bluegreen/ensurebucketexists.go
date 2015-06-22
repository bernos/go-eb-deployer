package bluegreen

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
)


func ensureBucketExists(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {

	var (
		bucket string
		exists bool
		err    error
	)

	log.Printf("Ensuring bucket exists...")

	s3Service := services.NewS3Service(s3.New(ctx.AwsConfig))

	if bucket, err = ctx.Bucket(); err != nil {
		return err
	}

	if exists, err = s3Service.BucketExists(bucket); err != nil {
		return err
	}

	if exists {
		log.Printf("Bucket %s already exists", bucket)
		return next()
	}

	log.Printf("Creating bucket %s", bucket)

	if err = s3Service.CreateBucket(bucket); err != nil {
		return err
	}

	log.Printf("Created bucket %s", bucket)

	return next()
}



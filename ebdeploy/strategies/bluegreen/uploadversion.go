package bluegreen

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
)

func uploadVersion(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {

	var (
		bucket        string
		err           error
		versionExists bool
	)

	s3Service := services.NewS3Service(s3.New(ctx.AwsConfig))
	ebService := services.NewEBService(elasticbeanstalk.New(ctx.AwsConfig))
	key := ctx.Version + ".zip"

	if versionExists, err = ebService.ApplicationVersionExists(ctx.Configuration.ApplicationName, ctx.Version); err != nil {
		return err
	} else if versionExists {
		return errors.New("Version " + ctx.Version + " already exists")
	}

	if bucket, err = ctx.Bucket(); err != nil {
		return err
	}

	if err = s3Service.UploadFile(ctx.SourceBundle, bucket, key); err != nil {
		return err
	}

	log.Printf("Uploaded version %s", ctx.Version)

	if err = ebService.CreateApplicationVersion(ctx.Configuration.ApplicationName, ctx.Version, bucket, key); err != nil {
		return err
	}

	log.Printf("Created version %s", ctx.Version)

	return next()
}



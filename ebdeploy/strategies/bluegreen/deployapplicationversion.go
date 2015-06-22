package bluegreen

import (

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
)

func deployApplicationVersion(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {

	log.Printf("Deploying version %s to environment %s", ctx.Version, ctx.TargetEnvironment.Name)

	client := elasticbeanstalk.New(ctx.AwsConfig)
	ebService := services.NewEBService(client)
	params := &elasticbeanstalk.CreateEnvironmentInput{
		ApplicationName:   aws.String(ctx.Configuration.ApplicationName),
		EnvironmentName:   aws.String(ctx.TargetEnvironment.Name),
		CNAMEPrefix:       aws.String(ctx.TargetEnvironment.CNAME),
		OptionSettings:    []*elasticbeanstalk.ConfigurationOptionSetting(ctx.Configuration.OptionSettings),
		Tags:              []*elasticbeanstalk.Tag(ctx.Configuration.Tags),
		SolutionStackName: aws.String(ctx.Configuration.SolutionStackName),
		Tier:              ctx.Configuration.Tier,
		VersionLabel:      aws.String(ctx.Version),
	}

	if resp, err := client.CreateEnvironment(params); err != nil {
		return err
	} else {

		log.Printf("Deployed: %v", resp)
		log.Printf("Waiting for environment health to go green")

		if err := ebService.WaitForEnvironment(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, func(e *elasticbeanstalk.EnvironmentDescription) bool {
			return *e.Health == "Green"
		}); err != nil {
			return err
		} else {
			log.Printf("Environment is green")
			return next()
		}
	}
}


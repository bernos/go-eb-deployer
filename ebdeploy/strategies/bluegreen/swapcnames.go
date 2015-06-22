package bluegreen

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
)

func swapCnames(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {
	log.Printf("Swapping cnames")

	client := elasticbeanstalk.New(ctx.AwsConfig)
	ebService := services.NewEBService(client)

	if environments, err := ebService.GetEnvironments(ctx.Configuration.ApplicationName); err == nil {

		activeCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, true)
		inactiveCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, false)

		activeEnvironment := findEnvironment(environments, cnamePredicate(activeCname))
		inactiveEnvironment := findEnvironment(environments, cnamePredicate(inactiveCname))

		if activeEnvironment != nil && inactiveEnvironment != nil {

			params := &elasticbeanstalk.SwapEnvironmentCNAMEsInput{
				DestinationEnvironmentID: aws.String(*activeEnvironment.EnvironmentID),
				SourceEnvironmentID:      aws.String(*inactiveEnvironment.EnvironmentID),
			}

			if _, err := client.SwapEnvironmentCNAMEs(params); err != nil {
				return err
			}
		} else if activeEnvironment == nil {
			return errors.New("No active environment to swap cname with")
		}

		log.Printf("Successfully swapped cnames")

		return next()
	} else {
		return err
	}
}


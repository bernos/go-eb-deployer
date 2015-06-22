package bluegreen

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"github.com/bernos/go-eb-deployer/ebdeploy/services"
	"log"
)

func prepareTargetEnvironment(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {

	ebService := services.NewEBService(elasticbeanstalk.New(ctx.AwsConfig))

	if environments, err := ebService.GetEnvironments(ctx.Configuration.ApplicationName); err == nil {

		log.Printf("Found %d existing environments", len(environments))
		requiresTerminate := false
		activeCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, true)
		inactiveCname := calculateCnamePrefix(ctx.Configuration.ApplicationName, ctx.Environment, false)

		activeEnvironment := findEnvironment(environments, cnamePredicate(activeCname))
		inactiveEnvironment := findEnvironment(environments, cnamePredicate(inactiveCname))

		if activeEnvironment != nil && inactiveEnvironment != nil {
			log.Println("Both active and inactive environments were found. Inactive environment will be terminated.")

			ctx.TargetEnvironment = &pipeline.TargetEnvironment{
				Name:     *inactiveEnvironment.EnvironmentName,
				CNAME:    inactiveCname,
				IsActive: false,
			}

			requiresTerminate = true
		} else if activeEnvironment == nil && inactiveEnvironment == nil {
			log.Println("Neither active nor inactive environments were found. Deploying directly to active environment")

			ctx.TargetEnvironment = &pipeline.TargetEnvironment{
				Name:     calculateEnvironmentName(ctx.Environment, "a"),
				CNAME:    activeCname,
				IsActive: true,
			}
		} else if activeEnvironment != nil {
			activeSuffix := getSuffixFromEnvironmentName(*activeEnvironment.EnvironmentName)
			inactiveSuffix := "a"

			if activeSuffix == "a" {
				inactiveSuffix = "b"
			}

			log.Printf("Active environment '%s' found. Deploying to inactive environment '%s'", activeSuffix, inactiveSuffix)

			ctx.TargetEnvironment = &pipeline.TargetEnvironment{
				Name:     calculateEnvironmentName(ctx.Environment, inactiveSuffix),
				CNAME:    inactiveCname,
				IsActive: false,
			}
		} else {
			return errors.New("Current environment state is not recognized. Please there is either one recognised, active environment, or no environments at all")
		}

		ctx.TargetEnvironment.Url = "http://" + ctx.TargetEnvironment.CNAME + ".elasticbeanstalk.com"

		done := make(chan struct{})
		defer close(done)
		ebService.LogEnvironmentEvents(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, done)

		if requiresTerminate {
			if err := ebService.TerminateEnvironment(*inactiveEnvironment.EnvironmentID); err != nil {
				return err
			}

			if err := ebService.WaitForEnvironment(ctx.Configuration.ApplicationName, ctx.TargetEnvironment.Name, func(e *elasticbeanstalk.EnvironmentDescription) bool {
				return *e.Status == "Terminated"
			}); err != nil {
				return err
			}
		}

		log.Printf("active environment: %v", activeEnvironment)
		log.Printf("inactive environment: %v", inactiveEnvironment)
		log.Printf("Target environment: %v", ctx.TargetEnvironment)

		return next()
	} else {
		return err
	}
}



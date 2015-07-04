package bluegreen

import (
	"log"

	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
)

func deployResources(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {
	if ctx.Configuration.Resources != nil {
		log.Printf("Deploying resources")

		return next()
	}

	log.Printf("No resources to deploy. Skipping...")

	return next()
}

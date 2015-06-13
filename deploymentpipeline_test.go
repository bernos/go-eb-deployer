package ebdeploy

import (
	"log"
	"testing"
)

func TestPipeline(t *testing.T) {

	pipeline := new(DeploymentPipeline)

	pipeline.AddStep(func(ctx *DeploymentContext, next Continue) error {
		log.Print("One")
		return next()
	})

	pipeline.AddStep(func(ctx *DeploymentContext, next Continue) error {
		log.Print("Two")
		return next()
	})

	pipeline.AddStep(func(ctx *DeploymentContext, next Continue) error {
		log.Print("Three")
		return next()
	})

	pipeline.Run(new(DeploymentContext))
}

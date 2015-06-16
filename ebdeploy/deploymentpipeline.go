package ebdeploy

import (
	"errors"
	"strings"
)

type DeploymentStep func(*DeploymentContext, Continue) error

type Continue func() error

type DeploymentPipeline struct {
	steps []DeploymentStep
}

func GetPipeline(strategy string) (*DeploymentPipeline, error) {
	switch {
	case strings.ToLower(strategy) == "blue-green":
		return NewBlueGreenStrategy(), nil
	}
	return nil, errors.New("Unknown deployment strategy " + strategy)
}

func (d *DeploymentPipeline) AddStep(step DeploymentStep) *DeploymentPipeline {
	d.steps = append(d.steps, step)
	return d
}

func (d *DeploymentPipeline) Run(ctx *DeploymentContext) error {
	pipeline := reduce(ctx, d.steps, func() error { return nil })
	return pipeline()
}

func reduce(ctx *DeploymentContext, items []DeploymentStep, acc Continue) func() error {
	if len(items) == 0 {
		return acc
	}

	step, items := items[len(items)-1], items[:len(items)-1]

	newAcc := func() error {
		return step(ctx, acc)
	}

	return reduce(ctx, items, newAcc)
}

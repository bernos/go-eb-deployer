package ebdeploy

import (
	"log"

	"github.com/bernos/go-eb-deployer/ebdeploy/config"
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	_ "github.com/bernos/go-eb-deployer/ebdeploy/strategies"
)

func Deploy(options Options) error {

	var (
		cfg     *config.Configuration
		context *pipeline.DeploymentContext
		pipe    *pipeline.DeploymentPipeline
		err     error
	)

	if cfg, err = config.LoadConfigFromFile(options.Config); err != nil {
		return err
	}

	if context, err = pipeline.NewDeploymentContext(cfg, options.Environment, options.Package, options.Version); err != nil {
		return err
	}

	if pipe, err = pipeline.GetPipeline(cfg.Strategy); err != nil {
		return err
	}

	log.Printf("Deploying version %s from package %s to environment %s", context.Version, context.SourceBundle, context.Environment)

	if err = pipe.Run(context); err != nil {
		return err
	}

	return nil
}

type Options struct {
	Version     string
	Environment string
	Package     string
	Config      string
}

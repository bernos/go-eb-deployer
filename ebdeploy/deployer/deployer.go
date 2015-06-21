package deployer

import (
	"log"

	"github.com/bernos/go-eb-deployer/ebdeploy"
	_ "github.com/bernos/go-eb-deployer/ebdeploy/strategies"
)

func Deploy(options Options) error {

	var (
		config  *ebdeploy.Configuration
		context *ebdeploy.DeploymentContext
		pipe    *ebdeploy.DeploymentPipeline
		err     error
	)

	if config, err = ebdeploy.LoadConfigFromFile(options.Config); err != nil {
		return err
	}

	if context, err = ebdeploy.NewDeploymentContext(config, options.Environment, options.Package, options.Version); err != nil {
		return err
	}

	if pipe, err = ebdeploy.GetPipeline(config.Strategy); err != nil {
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

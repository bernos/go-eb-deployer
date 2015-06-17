package ebdeploy

import (
	"log"
)

func Deploy(options Options) error {

	var (
		config   *Configuration
		context  *DeploymentContext
		pipeline *DeploymentPipeline
		err      error
	)

	if config, err = LoadConfigFromFile(options.Config); err != nil {
		return err
	}

	if context, err = NewDeploymentContext(config, options.Environment, options.Package, options.Version); err != nil {
		return err
	}

	if pipeline, err = GetPipeline(config.Strategy); err != nil {
		return err
	}

	log.Printf("Deploying version %s from package %s to environment %s", context.Version, context.SourceBundle, context.Environment)

	if err = pipeline.Run(context); err != nil {
		return err
	}

	return nil
}

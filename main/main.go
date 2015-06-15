package main

import (
	"flag"
	"fmt"
	"github.com/bernos/go-eb-deployer"
)

func main() {
	// environment, sourcebundle, version
	version := flag.String("version", "", "Version number")
	environment := flag.String("environment", "", "Environment to deploy to")
	sourceBundle := flag.String("package", "", "Package to deploy")
	configFile := flag.String("config", "", "Deployment config file")

	flag.Parse()

	if config, err := ebdeploy.LoadConfigFromFile(*configFile); err == nil {
		if context, err := ebdeploy.NewDeploymentContext(config, *environment, *sourceBundle, *version); err == nil {

			fmt.Println("Version:", context.Version)
			fmt.Println("Environment:", context.Environment)
			fmt.Println("SourceBundle:", context.SourceBundle)

			if pipeline, err := ebdeploy.GetPipeline(config.Strategy); err == nil {
				pipelineErr := pipeline.Run(context)

				if pipelineErr == nil {
					fmt.Println("Done")
				} else {
					panic(pipelineErr)
				}
			}

		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

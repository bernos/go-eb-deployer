package main

import (
	"flag"
	"github.com/bernos/go-eb-deployer/ebdeploy/deployer"
)

func ReadOptions() deployer.Options {
	options := deployer.Options{}

	flag.StringVar(&options.Version, "version", "", "Version number")
	flag.StringVar(&options.Environment, "environment", "", "Environment to deploy to")
	flag.StringVar(&options.Package, "package", "", "Package to deploy")
	flag.StringVar(&options.Config, "config", "", "Deployment config file")

	flag.Parse()

	return options
}

func main() {
	options := ReadOptions()

	if err := deployer.Deploy(options); err != nil {
		panic(err)
	}
}

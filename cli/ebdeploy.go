package main

import (
	"log"

	"github.com/bernos/go-eb-deployer/ebdeploy"
	"github.com/ttacon/chalk"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ReadOptions() ebdeploy.Options {
	options := ebdeploy.Options{}
	version := kingpin.Flag("version", "Version label.").Short('v').String()
	environment := kingpin.Flag("environment", "Environment to deploy to").Short('e').Required().String()
	bundle := kingpin.Flag("package", "Package to deploy").Short('p').Required().String()
	cfg := kingpin.Flag("config", "Deployment configuration json file").Short('c').Required().String()

	kingpin.Version("0.0.1")
	kingpin.Parse()

	options.Version = *version
	options.Environment = *environment
	options.Package = *bundle
	options.Config = *cfg

	return options
}

func main() {
	options := ReadOptions()

	if err := ebdeploy.Deploy(options); err != nil {
		log.Fatalf(chalk.Red.Color("ERROR: %s"), err)
	}
}

package ebdeploy

import (
	"testing"
)

var configJson = []byte(`{
	"ApplicationName" : "My application",
	"SolutionStackName" : "My solution stack",
	"Region" : "ap-southeast-2",
	"Bucket" : "my-bucket",
	"Tags" : [{
		"Key" : "Tag one",
		"Value" : "Value one"
	},
	{
		"Key" : "Tag Two",
		"Value" : "Value two"
	}],
	"OptionSettings" : [{
		"Namespace" : "a:b:c",
		"OptionName" : "optionone",
		"Value" : "valueone"
	}],
	"Tier" : {
		"Name" : "Web",
		"Type" : "WebType",
		"Version" : "1"
	},
	"Resources" : {
		"TemplateFile" : "my-resources.json",
		"Capabilities" : ["BLAH"],
		"Outputs" : [{
			"Name" : "OutputOne",
			"Namespace" : "d:e:f",
			"OptionName" : "blerg"
		}] 
	},
	"Environments" : [{
		"Name" : "Dev",
		"Description" : "Dev environment",
		"Tags" : [{
			"Key" : "One",
			"Value" : "Value for dev env"
		},
		{
			"Key" : "Tag Two",
			"Value" : "Value from dev"
		}],
		"OptionSettings" : [{
			"Namespace" : "a:b:c",
			"OptionName": "OptionOne",
			"Value":"value for dev env"
		}]
	},
	{
		"Name" : "Production",
		"Description" : "Production environment",
		"Tags" : [{
			"Key" : "One",
			"Value" : "Value for prod env"
		}],
		"OptionSettings" : [{
			"Namespace" : "a:b:c",
			"OptionName": "OptionOne",
			"Value":"value for prod env"
		}]
	}]
}`)

func loadConfig(t *testing.T) Configuration {
	config, err := LoadConfigFromJson(configJson)

	if err != nil {
		t.Errorf("Failed to parse config json %s", err)
	}

	return config
}

func TestHasEnvironment(t *testing.T) {
	config := loadConfig(t)

	if !config.HasEnvironment("Dev") {
		t.Errorf("Failed to find known environment in config")
	}

	if config.HasEnvironment("Blah") {
		t.Errorf("HasEnvironment returned true for non-existent environment")
	}
}

func TestNormalizeEnvironmentTags(t *testing.T) {
	config := loadConfig(t)

	if env, err := config.GetEnvironment("Dev"); err != nil {
		t.Errorf("%s", err)
	} else {
		tags := config.normalizeEnvironmentTags(env.Tags)

		if len(tags) != 3 {
			t.Errorf("Expected %d tags, but found %d", 3, len(tags))
		}
	}
}

func TestNormalizeEnvironment(t *testing.T) {
	config := loadConfig(t)

	env, _ := config.GetEnvironment("Dev")
	config.normalizeEnvironment(env)

	if len(env.Tags) != 3 {
		t.Errorf("Expected %d tags, but found %d", 3, len(env.Tags))
	}
}

func TestNormalizeEnvironments(t *testing.T) {
	config := loadConfig(t)

	config.normalizeEnvironments()

	env, _ := config.GetEnvironment("Dev")

	if len(env.Tags) != 3 {
		t.Errorf("Expected %d tags, but found %d", 3, len(env.Tags))
	}
}

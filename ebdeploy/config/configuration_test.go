package config

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
		"OptionName" : "OptionOne",
		"Value" : "valueone"
	},
	{
		"Namespace" : "e:y:w",
		"OptionName" : "optiontwo",
		"Value" : "foo"
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
		},
		{
			"Namespace" : "d:e:f",
			"OptionName" : "OptionOne",
			"Value":"value"
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

func loadConfig(t *testing.T) *Configuration {
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

func TestNormalize(t *testing.T) {
	config := loadConfig(t)

	env, _ := config.GetEnvironment("Dev")

	if len(env.Tags) != 3 {
		t.Errorf("Expected %d tags, but found %d", 3, len(env.Tags))
	}

	if len(env.OptionSettings) != 3 {
		t.Errorf("Expected %d option settings, but found %d", 3, len(env.OptionSettings))
	}

	if tag := env.Tags.GetTag("Tag Two"); tag != nil {
		expected := "Value from dev"
		actual := *tag.Value
		if actual != expected {
			t.Errorf("Expected tag value %s but found %s", expected, actual)
		}
	} else {
		t.Errorf("Failed to get tag")
	}

	if os := env.OptionSettings.GetOptionSetting("a:b:c", "OptionOne"); os != nil {
		expected := "value for dev env"
		actual := *os.Value
		if actual != expected {
			t.Errorf("Expected option setting value %s, but found %s", expected, actual)
		}
	} else {
		t.Errorf("Failed to find option setting")
	}
}

func TestNormalizeWithNoEnvironmentTags(t *testing.T) {
	cfg := []byte(`{
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
			"Description" : "Dev environment"
		}]
	}`)

	config, err := LoadConfigFromJson(cfg)

	if err != nil {
		t.Errorf("Failed to load from json %s", err)
	}

	if env, err := config.GetEnvironment("Dev"); err == nil {
		expected := 2
		actual := len(env.Tags)

		if expected != actual {
			t.Errorf("Expected %d tags but found %d", expected, actual)
		}
	} else {
		t.Errorf("Failed to get environment %s", err)
	}
}

func TestNormalizeWithNoConfigTags(t *testing.T) {
	cfg := []byte(`{
		"ApplicationName" : "My application",
		"SolutionStackName" : "My solution stack",
		"Region" : "ap-southeast-2",
		"Bucket" : "my-bucket",
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
			"Description" : "Dev environment"
		}]
	}`)

	config, err := LoadConfigFromJson(cfg)

	if err != nil {
		t.Errorf("Failed to load from json %s", err)
	}

	if env, err := config.GetEnvironment("Dev"); err == nil {
		expected := 0
		actual := len(env.Tags)

		if expected != actual {
			t.Errorf("Expected %d tags but found %d", expected, actual)
		}
	} else {
		t.Errorf("Failed to get environment %s", err)
	}
}

package ebdeploy

import (
	"encoding/json"
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

func TestUnmarshall(t *testing.T) {
	var config Configuration

	err := json.Unmarshal(configJson, &config)

	if err != nil {
		t.Errorf("Failed to unmarshall configuration json %s", err)
	}
}

func TestHasEnvironment(t *testing.T) {
	var config Configuration

	json.Unmarshal(configJson, &config)

	if !config.HasEnvironment("Dev") {
		t.Errorf("Failed to find known environment in config")
	}
}

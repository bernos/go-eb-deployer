package ebdeploy

import (
	"encoding/json"
	"fmt"
	//	"log"
)

type Tag struct {
	Key, Value string
}

type Tags []Tag

func (ts Tags) Contains(tag Tag) bool {
	for _, t := range ts {
		if t.Key == tag.Key {
			return true
		}
	}
	return false
}

type OptionSetting struct {
	Namespace, OptionName, Value string
}

type OptionSettings []OptionSetting

func (os OptionSettings) Contains(optionSetting OptionSetting) bool {
	for _, o := range os {
		if o.Namespace == optionSetting.namespace && o.OptionName == optionSetting.optionName {
			return true
		}
	}
	return false
}

type Tier struct {
	Name, Type, Version string
}

type Output struct {
	Name, Namespace, OptionName string
}

type Resources struct {
	TemplateFile string
	Capabilities []string
	Outputs      []Output
}

type Environment struct {
	Name, Description string
	Tags              []Tag
	OptionSettings    []OptionSetting
}

type Environments []Environment

type Configuration struct {
	ApplicationName   string
	SolutionStackName string
	Region            string
	Bucket            string
	Tags              Tags
	OptionSettings    OptionSettings
	Tier              Tier
	Resources         Resources
	Environments      Environments
}

func (c *Configuration) HasEnvironment(name string) bool {
	for _, env := range c.Environments {
		if env.Name == name {
			return true
		}
	}
	return false
}

func (c *Configuration) GetEnvironment(name string) (*Environment, error) {
	for _, env := range c.Environments {
		if env.Name == name {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("Environment %s not found", name)
}

func (c *Configuration) normalize() {
	c.normalizeEnvironments()
}

func (c *Configuration) normalizeEnvironments() {
	for i, _ := range c.Environments {
		env := &c.Environments[i]
		c.normalizeEnvironment(env)
	}
}

func (c *Configuration) normalizeEnvironment(environment *Environment) {
	environment.Tags = c.normalizeEnvironmentTags(environment.Tags)
}

func (c *Configuration) normalizeEnvironmentTags(environmentTags Tags) Tags {
	for _, t := range c.Tags {
		if !environmentTags.Contains(t) {
			environmentTags = append(environmentTags, t)
		}
	}
	return environmentTags
}

func LoadConfigFromJson(b []byte) (Configuration, error) {
	var config Configuration

	if err := json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	config.normalize()

	return config, nil
}

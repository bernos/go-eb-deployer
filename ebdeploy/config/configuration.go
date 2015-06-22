package config

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"io/ioutil"
)

type Tags []*elasticbeanstalk.Tag

func (ts Tags) Contains(tag *elasticbeanstalk.Tag) bool {
	for _, t := range ts {
		if *t.Key == *tag.Key {
			return true
		}
	}
	return false
}

func (ts Tags) GetTag(key string) *elasticbeanstalk.Tag {
	for i, t := range ts {
		if *t.Key == key {
			return ts[i]
		}
	}
	return nil
}

type OptionSettings []*elasticbeanstalk.ConfigurationOptionSetting

func (os OptionSettings) Contains(optionSetting *elasticbeanstalk.ConfigurationOptionSetting) bool {
	for _, o := range os {
		if *o.Namespace == *optionSetting.Namespace && *o.OptionName == *optionSetting.OptionName {
			return true
		}
	}
	return false
}

func (os OptionSettings) GetOptionSetting(namespace string, optionName string) *elasticbeanstalk.ConfigurationOptionSetting {
	for i, o := range os {
		if *o.Namespace == namespace && *o.OptionName == optionName {
			return os[i]
		}
	}
	return nil
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
	Tags              Tags
	OptionSettings    OptionSettings
}

func (e *Environment) MergeCommonTags(tags Tags) {
	for _, t := range tags {
		if !e.Tags.Contains(t) {
			e.Tags = append(e.Tags, t)
		}
	}
}

func (e *Environment) MergeCommonOptionSettings(optionSettings OptionSettings) {
	for _, os := range optionSettings {
		if !e.OptionSettings.Contains(os) {
			e.OptionSettings = append(e.OptionSettings, os)
		}
	}
}

type Environments []Environment

type Configuration struct {
	ApplicationName   string
	SolutionStackName string
	Strategy          string
	Region            string
	Bucket            string
	Tags              Tags
	OptionSettings    OptionSettings
	Tier              *elasticbeanstalk.EnvironmentTier
	Resources         Resources
	Environments      Environments
	SmokeTestUrl      string
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
	environment.MergeCommonTags(c.Tags)
	environment.MergeCommonOptionSettings(c.OptionSettings)
}

func LoadConfigFromJson(b []byte) (*Configuration, error) {
	config := new(Configuration)

	if err := json.Unmarshal(b, config); err != nil {
		return nil, err
	}

	config.normalize()

	return config, nil
}

func LoadConfigFromFile(file string) (*Configuration, error) {
	if bytes, err := ioutil.ReadFile(file); err == nil {
		return LoadConfigFromJson(bytes)
	} else {
		return nil, err
	}
}

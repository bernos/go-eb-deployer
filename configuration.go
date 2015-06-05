package ebdeploy

type Tag struct {
	Key, Value string
}

type OptionSetting struct {
	Namespace, OptionName, Value string
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

type Configuration struct {
	ApplicationName   string
	SolutionStackName string
	Region            string
	Bucket            string
	Tags              []Tag
	OptionSettings    []OptionSetting
	Tier              Tier
	Resources         Resources
	Environments      []Environment
}

func (c *Configuration) HasEnvironment(name string) bool {
	for _, env := range c.Environments {
		if env.Name == name {
			return true
		}
	}
	return false
}

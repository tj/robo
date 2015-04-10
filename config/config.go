package config

import "github.com/tj/robo/task"
import "gopkg.in/yaml.v2"
import "io/ioutil"

// Config.
type Config struct {
	Tasks     map[string]*task.Task `yaml:",inline"`
	Templates struct {
		List string
		Help string
	}
}

// New configuration loaded from `file`.
func New(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewString(string(b))
}

// NewString configuration from string.
func NewString(s string) (*Config, error) {
	c := new(Config)

	err := yaml.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}

	for name, task := range c.Tasks {
		task.Name = name
	}

	return c, nil
}

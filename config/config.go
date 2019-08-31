package config

import (
	"bytes"
	"io/ioutil"
	"os/user"
	"path"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/tj/robo/task"
)

// Config represents the main YAML configuration
// loaded for Robo tasks.
type Config struct {
	File      string
	Tasks     map[string]*task.Task `yaml:",inline"`
	Variables map[string]interface{}
	Templates struct {
		List      string
		Help      string
		Variables string
	}
}

// Eval evaluates the config by interpolating
// all templates using the variables.
func (c *Config) Eval() error {
	for _, task := range c.Tasks {
		err := interpolate(
			c.Variables,
			&task.Command,
			&task.Summary,
			&task.Script,
			&task.Exec,
		)
		if err != nil {
			return err
		}

		for i, item := range task.Env {
			if err := interpolate(c.Variables, &item); err != nil {
				return err
			}
			task.Env[i] = item
		}
	}
	return nil
}

// New configuration loaded from `file`.
func New(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c, err := NewString(string(b))
	if err != nil {
		return nil, err
	}
	c.File = file

	// Initialize variables if needed.
	if c.Variables == nil {
		c.Variables = make(map[string]interface{})
	}

	// Expose robo's internal variables
	// but respect users who override them.
	if _, ok := c.Variables["robo"]; !ok {
		c.Variables["robo"] = map[string]string{
			"path": path.Dir(c.File),
			"file": c.File,
		}
	}

	// Add the current user.
	if _, ok := c.Variables["user"]; !ok {
		if user, err := user.Current(); err == nil {
			c.Variables["user"] = map[string]string{
				"name":     user.Name,
				"username": user.Username,
				"home":     user.HomeDir,
			}
		}
	}

	// Interpolate variables.
	if err := c.Eval(); err != nil {
		return nil, err
	}

	return c, nil
}

// NewString configuration from string.
func NewString(s string) (*Config, error) {
	c := new(Config)

	// unmarshal
	err := yaml.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}

	// assign .Name
	for name, task := range c.Tasks {
		task.Name = name
	}

	return c, nil
}

// Apply interpolation against the given strings.
func interpolate(v interface{}, s ...*string) error {
	for _, p := range s {
		ret, err := eval(*p, v)
		if err != nil {
			return err
		}
		*p = ret
	}
	return nil
}

// Evaluate template against `v`.
func eval(s string, v interface{}) (string, error) {
	t, err := template.New("task").Parse(s)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.Execute(&b, v)
	if err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}

package config

import "github.com/tj/robo/task"
import "gopkg.in/yaml.v2"
import "text/template"
import "io/ioutil"
import "bytes"

// Config.
type Config struct {
	Tasks     map[string]*task.Task `yaml:",inline"`
	Variables map[string]interface{}
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

	// unmarshal
	err := yaml.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}

	// assign .Name
	for name, task := range c.Tasks {
		task.Name = name
	}

	// interpolation
	for _, task := range c.Tasks {
		err := interpolate(c.Variables, &task.Command, &task.Summary, &task.Script, &task.Exec)
		if err != nil {
			return nil, err
		}
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

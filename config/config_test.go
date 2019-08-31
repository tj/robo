package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/tj/robo/config"
)

var s = `
foo:
  summary: Command foo.
  command: echo "foo"

bar:
  summary: Command bar.
  command: echo "bar"

stage:
  summary: Run commands against stage.
  command: ssh {{.hosts.stage}} -t robo

prod:
  summary: Run commands against prod.
  command: ssh {{.hosts.prod}} -t robo
  env: ["H={{.hosts.prod}}"]

templates:
  list: testing

variables:
  hosts:
    prod: bastion-prod
    stage: bastion-stage
`

func TestNewString(t *testing.T) {
	c, err := config.NewString(s)
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(c.Tasks))

	assert.Equal(t, nil, c.Eval())

	assert.Equal(t, `ssh bastion-stage -t robo`, c.Tasks["stage"].Command)
	assert.Equal(t, `ssh bastion-prod -t robo`, c.Tasks["prod"].Command)

	assert.Equal(t, `foo`, c.Tasks["foo"].Name)
	assert.Equal(t, `Command foo.`, c.Tasks["foo"].Summary)
	assert.Equal(t, `echo "foo"`, c.Tasks["foo"].Command)

	assert.Equal(t, `Command bar.`, c.Tasks["bar"].Summary)
	assert.Equal(t, `echo "bar"`, c.Tasks["bar"].Command)

	assert.Equal(t, []string{"H=bastion-prod"}, c.Tasks["prod"].Env)

	assert.Equal(t, `testing`, c.Templates.List)
}

func TestNew(t *testing.T) {
	b := []byte(s)
	f, err := ioutil.TempFile("", "")
	assert.Equal(t, nil, err)

	file := f.Name()
	defer os.Remove(file)

	_, err = f.Write(b)
	assert.Equal(t, nil, err)

	f.Close()

	c, err := config.New(file)
	assert.Equal(t, nil, err)
	assert.Equal(t, file, c.File)
}

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
  region:
    euw: europe-west
  hosts:
    stage: bastion-stage
    prod: bastion-prod
  dns: "{{ .region.euw }}.{{ .hosts.prod }}"
  command: "$(true && echo $?)"
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

	// test variables section interpolation
	assert.Equal(t, "europe-west.bastion-prod", c.Variables["dns"])
	assert.Equal(t, 0, c.Variables["command"])

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
package config_test

import "github.com/bmizerany/assert"
import "github.com/tj/robo/config"
import "testing"

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

	assert.Equal(t, `ssh bastion-stage -t robo`, c.Tasks["stage"].Command)
	assert.Equal(t, `ssh bastion-prod -t robo`, c.Tasks["prod"].Command)

	assert.Equal(t, `foo`, c.Tasks["foo"].Name)
	assert.Equal(t, `Command foo.`, c.Tasks["foo"].Summary)
	assert.Equal(t, `echo "foo"`, c.Tasks["foo"].Command)

	assert.Equal(t, `Command bar.`, c.Tasks["bar"].Summary)
	assert.Equal(t, `echo "bar"`, c.Tasks["bar"].Command)

	assert.Equal(t, `testing`, c.Templates.List)
}

package interpolation

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/tj/robo/task"
)

func TestVars_whenValueReferencesOtherKey_shouldReplaceAccordingly(t *testing.T) {
	vars := map[string]interface{}{
		"foo": "Hello",
		"bar": "{{ .foo }} World!",
	}
	err := Vars(&vars)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Hello", vars["foo"])
	assert.Equal(t, "Hello World!", vars["bar"])
}

func TestVars_whenValueIsCommand_shouldReplaceWithCommandResult(t *testing.T) {
	vars := map[string]interface{}{
		"foo": "$(echo Hello)",
		"bar": map[string]interface{}{
			"sub": "$(echo World!)",
		},
	}

	err := Vars(&vars)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Hello", vars["foo"])
	assert.Equal(t, "World!", vars["bar"].(map[interface{}]interface{})["sub"])
}

func TestTasks(t *testing.T) {
	tk := task.Task{
		Summary: "This task handles {{ .foo }} World!",
		Before: []*task.Runnable{
			{Command: "{{ .foo }}"},
			{Script: "{{ .foo }}"},
			{Exec: "{{ .foo }}"},
		},
		After: []*task.Runnable{
			{Command: "{{ .bar }}"},
			{Script: "{{ .bar }}"},
			{Exec: "{{ .bar }}"},
		},
		Command: "echo {{ .foo }} World!",
		Script:  "/path/to/{{ .foo }}.sh",
		Exec:    "{{ .foo }} World!",
		Env:     []string{"bar={{ .foo }} World!"},
		Usage:   "robo {{.bar}}",
		Examples: []*task.Example{
			{Description: "{{ .foo }} Example!", Command: "robo {{.bar}}"},
		},
	}

	vars := map[string]interface{}{"foo": "Hello", "bar": "Bye"}

	err := Tasks(map[string]*task.Task{"tk": &tk}, vars)
	assert.Equal(t, nil, err)
	assert.Equal(t, "This task handles Hello World!", tk.Summary)
	assert.Equal(t, "echo Hello World!", tk.Command)
	assert.Equal(t, "/path/to/Hello.sh", tk.Script)
	assert.Equal(t, "Hello World!", tk.Exec)
	assert.Equal(t, "Hello", tk.Before[0].Command)
	assert.Equal(t, "Hello", tk.Before[1].Script)
	assert.Equal(t, "Hello", tk.Before[2].Exec)
	assert.Equal(t, "Bye", tk.After[0].Command)
	assert.Equal(t, "Bye", tk.After[1].Script)
	assert.Equal(t, "Bye", tk.After[2].Exec)
	assert.Equal(t, "bar=Hello World!", tk.Env[0])
	assert.Equal(t, "robo Bye", tk.Usage)
	assert.Equal(t, "Hello Example!", tk.Examples[0].Description)
	assert.Equal(t, "robo Bye", tk.Examples[0].Command)
}

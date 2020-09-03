package interpolation

import (
	"github.com/bmizerany/assert"
	"github.com/tj/robo/task"
	"testing"
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
		"bar":
			map[string]interface{}{
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
		Summary:    "This task handles {{ .foo }} World!",
		Command:    "echo {{ .foo }} World!",
		Script:     "/path/to/{{ .foo }}.sh",
		Exec:       "{{ .foo }} World!",
		Env:        []string{"bar={{ .foo }} World!"},
	}

	vars := map[string]interface{}{"foo": "Hello"}

	err := Tasks(map[string]*task.Task{"tk": &tk}, vars)
	assert.Equal(t, nil, err)
	assert.Equal(t, "This task handles Hello World!", tk.Summary)
	assert.Equal(t, "echo Hello World!", tk.Command)
	assert.Equal(t, "/path/to/Hello.sh", tk.Script)
	assert.Equal(t, "Hello World!", tk.Exec)
	assert.Equal(t, "bar=Hello World!", tk.Env[0])
}

package cli

import (
	"github.com/bmizerany/assert"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	vars := map[string]interface{}{
		"root": "foo",
		"one":
		map[interface{}]interface{}{
			"two": "bar",
			"three": map[interface{}]interface{}{
				"hello": "world",
			},
		},
	}

	flattened := flatten("", reflect.ValueOf(vars))

	assert.Equal(t, 3, len(flattened))
	assert.Equal(t, "foo", flattened[".root"])
	assert.Equal(t, "bar", flattened[".one.two"])
	assert.Equal(t, "world", flattened[".one.three.hello"])
}
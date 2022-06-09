// Package interpolation provides methods to interpolate user defined variables and robo tasks.
package interpolation

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/tj/robo/task"
	"gopkg.in/yaml.v2"
)

var commandPattern = regexp.MustCompile("\\$\\((.+)\\)")

// Vars interpolates a given map of interfaces (strings or submaps) with itself
// returning it mit populated template values.
func Vars(vars *map[string]interface{}) error {
	b, err := yaml.Marshal(*vars)
	if err != nil {
		return err
	}
	s := string(b)

	err = interpolate("variables", *vars, &s)
	if err != nil {
		return err
	}

	err = interpolateVariableCommands(&s)
	if err != nil {
		return fmt.Errorf("failed replacing variable placeholder with command result")
	}

	err = yaml.Unmarshal([]byte(s), vars)
	if err != nil {
		return err
	}
	return err
}

func interpolateVariableCommands(s *string) error {
	// find all commands
	matches := commandPattern.FindAllStringSubmatch(*s, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		cmdOut, err := captureCommandOutput(match[1])
		if err != nil {
			return fmt.Errorf("error while executing command. Error: %s", err)
		}
		*s = strings.ReplaceAll(*s, match[0], cmdOut)
	}
	return nil
}

// captureCommandOutput executes a command and captures the output which usually gets prompted to stdout.
func captureCommandOutput(args string) (string, error) {
	var cmd *exec.Cmd
	// try to use the user's default shell. If it is not set via env var fall back to `sh`.
	if defaultShell, ok := os.LookupEnv("SHELL"); ok {
		cmd = exec.Command(defaultShell, "-c", args)
	} else {
		cmd = exec.Command("sh", "-c", args)
	}
	var b bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSuffix(string(b.Bytes()), "\n"), err
}

// Tasks interpolates a given task with a set of data replacing placeholders
// in the command, summary, script, exec and envs property.
func Tasks(tasks map[string]*task.Task, data map[string]interface{}) error {
	for _, task := range tasks {
		// interpolate the tasks main fields
		err := interpolate(
			"task",
			data,
			&task.Command,
			&task.Summary,
			&task.Script,
			&task.Exec,
		)
		if err != nil {
			return err
		}

		// interpolate a task's environment data
		for i, item := range task.Env {
			if err := interpolate("env-var", data, &item); err != nil {
				return err
			}
			task.Env[i] = item
		}

		// interpolate a task's before and after steps
		if err := interpolateOptionals("before", task.Before, data); err != nil {
			return err
		}
		if err := interpolateOptionals("after", task.After, data); err != nil {
			return err
		}
	}
	return nil
}

func interpolateOptionals(id string, rs []*task.Runnable, data map[string]interface{}) error {
	for i, step := range rs {
		err := interpolate(
			id,
			data,
			&step.Command,
			&step.Exec,
			&step.Script,
		)
		if err != nil {
			return err
		}
		rs[i] = step
	}
	return nil
}

// interpolate populates a given slice of templates with actual values provided
// in the data parameter.
func interpolate(name string, data interface{}, temps ...*string) error {
	for _, temp := range temps {
		t, err := template.New(name).Parse(*temp)
		if err != nil {
			return err
		}

		var b bytes.Buffer
		err = t.Execute(&b, data)
		if err != nil {
			return err
		}
		*temp = string(b.Bytes())
	}
	return nil
}

package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
)

// Example usage.
type Example struct {
	Description string
	Command     string
}

// Task definition.
type Task struct {
	LookupPath string
	Name       string `yaml:"-"`
	Summary    string
	Command    string
	Script     string
	Exec       string
	Usage      string
	Examples   []*Example
	Env        []string
	Before     []*Runnable
	After      []*Runnable
}

// Run the task and its preceding and succeding steps with `args`.
// - A failing before step will still allow the main task and the after steps to be executed
// - A failing task will always allow the after steps to be executed
func (t *Task) Run(args []string) []error {
	var errs []error
	if err := t.runOptionals("before", t.Before, args); err != nil {
		errs = append(errs, err)
	}

	// wrap command, script, exec into a runnable
	r := Runnable{Command: t.Command, Script: t.Script, Exec: t.Exec}

	if err := r.Run(t.LookupPath, args, t.Env); err != nil {
		errs = append(errs, fmt.Errorf("task '%s' failed. Error: %+v", t.Name, err))
	}

	if err := t.runOptionals("after", t.After, args); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (t *Task) runOptionals(id string, rs []*Runnable, args []string) error {
	for i, r := range rs {
		if err := r.Run(t.LookupPath, args, t.Env); err != nil {
			return fmt.Errorf("%s step #%d of task '%s' failed. Error: %+v", id, i+1, t.Name, err)
		}
	}
	return nil
}

type Runnable struct {
	Command string
	Script  string
	Exec    string
}

func (r *Runnable) Run(lookupPath string, args []string, env []string) error {
	if r.Exec != "" {
		return r.RunExec(args, env)
	}

	if r.Script != "" {
		return r.RunScript(lookupPath, args, env)
	}

	if r.Command != "" {
		return r.RunCommand(args, env)
	}

	return fmt.Errorf("nothing to run (add script, command, or exec key)")
}

// RunScript runs the target shell `script` file.
func (r *Runnable) RunScript(lookupPath string, args []string, env []string) error {
	var path = r.Script
	var bin = path

	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(lookupPath, r.Script)
		bin = path
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() && stat.Mode()&0100 == 0 {
		args = append([]string{path}, args...)
		bin = "sh"
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommand runs the `command` via the shell.
func (r *Runnable) RunCommand(args []string, env []string) error {
	args = append([]string{"-c", r.Command, "sh"}, args...)
	cmd := exec.Command("sh", args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunExec runs the `exec` command.
func (r *Runnable) RunExec(args []string, env []string) error {
	fields, err := shellwords.Parse(r.Exec)
	if err != nil {
		return err
	}

	bin := fields[0]
	path, err := exec.LookPath(bin)
	if err != nil {
		return err
	}

	envs := merge(os.Environ(), env)
	args = append(fields, args...)
	return syscall.Exec(path, args, envs)
}

// Merge merges the given two lists of env vars.
func merge(a, b []string) []string {
	var items = make(map[string]string)
	var ret []string

	for _, item := range a {
		if i := strings.Index(item, "="); i != -1 {
			key := item[:i]
			items[key] = item[i+1:]
		}
	}

	for _, item := range b {
		if i := strings.Index(item, "="); i != -1 {
			key := item[:i]
			items[key] = item[i+1:]
		}
	}

	for k, v := range items {
		ret = append(ret, k+"="+v)
	}

	return ret
}

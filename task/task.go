package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
}

// Run the task with `args`.
func (t *Task) Run(args []string) error {
	if t.Exec != "" {
		return t.RunExec(args)
	}

	if t.Script != "" {
		return t.RunScript(args)
	}

	if t.Command != "" {
		return t.RunCommand(args)
	}

	return fmt.Errorf("nothing to run (add script, command, or exec key)")
}

// RunScript runs the target shell `script` file.
func (t *Task) RunScript(args []string) error {
	var path = t.Script
	var bin = path

	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(t.LookupPath, t.Script)
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
	cmd.Env = append(os.Environ(), t.Env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommand runs the `command` via the shell.
func (t *Task) RunCommand(args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		commands := strings.Join(strings.Split(strings.TrimSpace(t.Command), "\n"), "&")
		commands += "/C" + commands
		cmd = exec.Command("cmd", commands)
	} else {
		args = append([]string{"-c", t.Command, "sh"}, args...)
		cmd = exec.Command("sh", args...)
	}
	cmd.Env = append(os.Environ(), t.Env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunExec runs the `exec` command.
func (t *Task) RunExec(args []string) error {
	fields, err := shellwords.Parse(t.Exec)
	if err != nil {
		return err
	}

	bin := fields[0]
	path, err := exec.LookPath(bin)
	if err != nil {
		return err
	}

	env := merge(os.Environ(), t.Env)
	args = append(fields, args...)
	return syscall.Exec(path, args, env)
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

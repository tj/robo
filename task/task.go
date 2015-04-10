package task

import "os/exec"
import "strings"
import "syscall"
import "fmt"
import "os"

// Example usage.
type Example struct {
	Description string
	Command     string
}

// Task definition.
type Task struct {
	Name     string `yaml:"-"`
	Summary  string
	Command  string
	Script   string
	Exec     string
	Usage    string
	Examples []*Example
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
	args = append([]string{t.Script}, args...)
	cmd := exec.Command("sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommand runs the `command` via the shell.
func (t *Task) RunCommand(args []string) error {
	args = append([]string{t.Command}, args...)
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunExec runs the `exec` command.
func (t *Task) RunExec(args []string) error {
	fields := strings.Fields(t.Exec)
	bin := fields[0]

	path, err := exec.LookPath(bin)
	if err != nil {
		return err
	}

	args = append(fields, args...)
	return syscall.Exec(path, args, os.Environ())
}

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/fatih/color"
	"github.com/tj/robo/config"
	"github.com/tj/robo/task"
)

// Template helpers.
var helpers = template.FuncMap{
	"magenta": color.MagentaString,
	"yellow":  color.YellowString,
	"green":   color.GreenString,
	"black":   color.BlackString,
	"white":   color.WhiteString,
	"blue":    color.BlueString,
	"cyan":    color.CyanString,
	"red":     color.RedString,
}

// List template.
var list = `
{{range .Tasks}}  {{cyan .Name}} – {{.Summary}}
{{end}}
`

// Quiet list.
var quiet = "{{range .Tasks}}{{ .Name }}\n{{end}}"

// Variables template.
var variables = `
{{- range $k, $v := . }}
{{cyan "%s" $k }}: {{$v}}
{{- end }}
`

// Help template.
var help = `
  {{cyan "Usage:"}}

    {{.Name}} {{.Usage}}

  {{cyan "Description:"}}

    {{.Summary}}
{{with .Examples}}
  {{cyan "Examples:"}}
  {{range .}}
    {{.Description}}
    $ {{.Command}}
  {{end}}{{end}}
`

// ListVariables outputs the variables defined.
func ListVariables(c *config.Config) {
	tmpl := t(variables)

	if c.Templates.Variables != "" {
		tmpl = t(c.Templates.Variables)
	}

	flattened := flatten("", reflect.ValueOf(c.Variables))
	tmpl.Execute(os.Stdout, flattened)
}

// List outputs the tasks defined.
func List(c *config.Config) {
	tmpl := t(list)

	if c.Templates.List != "" {
		tmpl = t(c.Templates.List)
	}

	tmpl.Execute(os.Stdout, c)
}

// ListNames lists task names.
func ListNames(c *config.Config) {
	tmpl := t(quiet)
	tmpl.Execute(os.Stdout, c)
}

// Help outputs the task help.
func Help(c *config.Config, name string) {
	task, ok := c.Tasks[name]

	if !ok {
		Fatalf("undefined task %q", name)
	}

	tmpl := t(help)

	if c.Templates.Help != "" {
		tmpl = t(c.Templates.Help)
	}

	tmpl.Execute(os.Stdout, task)
}

// Run the task.
func Run(c *config.Config, name string, args []string) {
	t, ok := c.Tasks[name]
	if !ok {
		Fatalf("undefined task %q", name)
	}
	lookupPath := filepath.Dir(c.File)
	t.LookupPath = lookupPath

	var errs []error
	if err := task.RunOptionals("before", "GLOBAL", c.Before, args, lookupPath, nil); err != nil {
		errs = append(errs, err)
	}

	if runErrs := t.Run(args); len(runErrs) > 0 {
		errs = append(errs, runErrs...)
	}

	if err := task.RunOptionals("after", "GLOBAL", c.After, args, lookupPath, nil); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		var msg string
		for _, err := range errs {
			msg += fmt.Sprintf("    - %+v\n", err)
		}
		Fatalf("error(s): \n%s", msg)
	}
}

// Fatalf writes to stderr and exits.
func Fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\n  %s\n\n", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Template helper.
func t(s string) *template.Template {
	return template.Must(template.New("").Funcs(helpers).Parse(s))
}

// flatten reduces a given map into a flattened map of strings having the path to a variable as a key
// and the actual value as a value. Resulting in ".path.to.key: value"
func flatten(key string, v reflect.Value) map[string]string {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	m := map[string]string{}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			for k, v := range flatten(key+"."+fmt.Sprintf("%v", k), v.MapIndex(k)) {
				m[k] = v
			}
		}
	default:
		m[key] = v.String()
	}
	return m
}

package cli

import "github.com/tj/robo/config"
import "github.com/fatih/color"
import "text/template"
import "fmt"
import "os"

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

// List outputs the tasks defined.
func List(c *config.Config) {
	tmpl := t(list)

	if c.Templates.List != "" {
		tmpl = t(c.Templates.List)
	}

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
	task, ok := c.Tasks[name]

	if !ok {
		Fatalf("undefined task %q", name)
	}

	err := task.Run(args)
	if err != nil {
		Fatalf("error: %s", err)
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

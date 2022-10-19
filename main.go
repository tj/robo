package main

import (
	"path/filepath"

	"github.com/tj/docopt"
	"github.com/tj/robo/cli"
	"github.com/tj/robo/config"
)

var version = "0.8.0"

const usage = `
  Usage:
    robo [-q] [--config file]
    robo <task> [<arg>...] [--config file]
    robo help [<task>] [--config file]
    robo variables [--config file]
    robo -h | --help
    robo --version

  Options:
    -c, --config file   config file to load [default: robo.yml]
    -h, --help          output help information
    -v, --version       output version
    -q, --quiet         output task names only

  Examples:

    output tasks
    $ robo

    output task help
    $ robo help mytask

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		cli.Fatalf("error parsing arguments: %s", err)
	}

	abs, err := filepath.Abs(args["--config"].(string))
	if err != nil {
		cli.Fatalf("cannot resolve --config: %s", err)
	}

	c, err := config.New(abs)
	if err != nil {
		cli.Fatalf("error loading configuration: %s", err)
	}

	switch {
	case args["help"].(bool):
		if name, ok := args["<task>"].(string); ok {
			cli.Help(c, name)
		} else {
			cli.List(c)
		}
	case args["variables"].(bool):
		cli.ListVariables(c)
	default:
		if name, ok := args["<task>"].(string); ok {
			cli.Run(c, name, args["<arg>"].([]string))
			return
		}

		if args["--quiet"].(bool) {
			cli.ListNames(c)
		} else {
			cli.List(c)
		}
	}
}

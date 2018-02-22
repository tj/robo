package main

import (
	"os"

	"github.com/tj/docopt"
	"github.com/tj/robo/cli"
	"github.com/tj/robo/config"
)

var version = "0.4.1"

const usage = `
  Usage:
    robo [--config file]
    robo <task> [<arg>...] [--config file]
    robo help [<task>] [--config file]
    robo variables [--config file]
    robo -h | --help
    robo --version

  Options:
    -c, --config file   config file to load
    -h, --help          output help information
    -v, --version       output version

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

	c := args["--config"]
	if c == nil {
		c = os.Getenv("ROBO_CONFIG")
		if c == nil || c == "" {
			cli.Fatalf("robo requires a config file passed via --config or the ROBO_CONFIG env var")
		}
	}

	file := c.(string)
	conf, err := config.New(file)
	if err != nil {
		cli.Fatalf("error loading configuration: %s", err)
	}

	switch {
	case args["help"].(bool):
		if name, ok := args["<task>"].(string); ok {
			cli.Help(conf, name)
		} else {
			cli.List(conf)
		}
	case args["variables"].(bool):
		cli.ListVariables(conf)
	default:
		if name, ok := args["<task>"].(string); ok {
			cli.Run(conf, name, args["<arg>"].([]string))
		} else {
			cli.List(conf)
		}
	}
}

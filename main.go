package main

import "github.com/tj/robo/config"
import "github.com/tj/robo/cli"
import "github.com/tj/docopt"

var Version = "0.1.0"

const Usage = `
  Usage:
    robo [--config file]
    robo <task> [<arg>...] [--config file]
    robo help [<task>] [--config file]
    robo -h | --help
    robo --version

  Options:
    -c, --config file   config file to load [default: robo.yml]
    -h, --help          output help information
    -v, --version       output version

  Examples:

    output tasks
    $ robo

    output task help
    $ robo help mytask

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, true)
	if err != nil {
		cli.Fatalf("error parsing arguments: %s", err)
	}

	file := args["--config"].(string)
	c, err := config.New(file)
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
	default:
		if name, ok := args["<task>"].(string); ok {
			cli.Run(c, name, args["<arg>"].([]string))
		} else {
			cli.List(c)
		}
	}
}

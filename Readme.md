
# Robo

 Simple YAML-based task runner written in Go.

 ![](http://img1.wikia.nocookie.net/__cb20130827204332/callofduty/images/f/fa/Giant_mech_Origins_BOII.png)

## Features

 - Not super slow
 - Not super obscure
 - No dependencies
 - Simple
 - That's it

## Installation

 Via `go-get`:

```
$ go get github.com/tj/robo
```

## Usage

 Command-line usage.

### Listing tasks

 Output tasks:

```
$ robo

  aws – amazon web services cli
  circle.open – open the repo in circle ci
  events – send data to the "events" topic
  push – push image from the current directory

```

### Task help

 Output task help:

```
$ robo help events

  Usage:

    events [project-id] [rate]

  Description:

    send data to the "events" topic

  Examples:

    Send 25 events a second to gy2d
    $ robo events gy2d 25

```

### Running tasks

 Regardless of task type (shell, exec, script) any additional arguments
 or given will be passed.

 For example suppose you have the following task:

```yml
aws:
  command: ssh tools aws
```

 You may then interact with the AWS cli as you would normally:

```
$ robo aws help
$ robo aws ec2 describe-instances
```

## Configuration

 Task configuration.

### Commands

 The most basic task simply runs a shell command:

```yml
hello:
  summary: some task
  command: echo world
```

 You may also define multi-line commands with YAML's `|`:

```yml
hello:
  summary: some task
  command: |
    echo hello
    echo world
```

### Exec

 The exec alternative lets you replace the robo image without
 the fork & shell, however note that shell features are not
 available (pipes, redirection, etc).

```yml
hello:
  summary: some task
  exec: echo hello
```

### Scripts

 Shell scripts may be used instead of inline commands:

```yml
hello:
  summary: some task
  script: path/to/script.sh
```

### Usage

 Tasks may optionally specify usage parameters, which display
 upon help output:

```yml
events:
  summary: send data to the "events" topic
  command: docker run -it events
  usage: "[project-id] [rate]"
```

### Examples

 Tasks may optionally specify any number of example commands, which
 display upon help output:

```yml
events:
  summary: send data to the "events" topic
  command: docker run -it events
  usage: "[project-id] [rate]"
  examples:
    - description: Send 25 events a second to gy2d
      command: robo events gy2d 25
```

### Templates

 Task `list` and `help` output may be re-configured, for example if you
 prefer to view usage information instead of the summary:

```yml
templates:
  list: |
    {{range .Tasks}}  {{cyan .Name}} – {{.Usage}}
    {{end}}
```

 Or perhaps something more verbose:

```yml
templates:
  list: |
    {{range .Tasks}}
      name: {{cyan .Name}}
      summary: {{.Summary}}
      usage: {{.Usage}}
    {{end}}
```

## Global tasks

 By default `./robo.yml` is loaded, however if you want global tasks
 you can simply alias to something like:

```
alias segment='robo --config ~/.robo.yml'
```

## Robo chaining

 You can easily use Robo to chain Robo, which is useful
 for multi-environment setups. For example:

```yml
prod:
  summary: production tasks
  exec: robo --config production.yml

stage:
  summary: stage tasks
  exec: robo --config stage.yml
```

 Or on remote boxes:

```yml
prod:
  summary: production tasks
  exec: ssh prod-tools -t robo --config production.yml

stage:
  summary: stage tasks
  exec: ssh stage-tools -t robo --config production.yml
```

## Why?

 We generally use Makefiles for project specific tasks, however
 the discoverability of global tasks within a large team is
 difficult unless there's good support for self-documentation,
 which Make is bad at.

 I'm aware of the million other solutions (Sake, Thor, etc) but
 I just wanted something fast without dependencies.

# License

 MIT
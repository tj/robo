
# Robo

 [![Build Status](https://travis-ci.org/tj/robo.svg?branch=master)](https://travis-ci.org/tj/robo)

 Simple YAML-based task runner written in Go.

 ![](http://img1.wikia.nocookie.net/__cb20130827204332/callofduty/images/f/fa/Giant_mech_Origins_BOII.png)

## Features

 - Not super slow
 - Not super obscure
 - No dependencies
 - Variables
 - Simple

## Installation

From [gobinaries.com](https://gobinaries.com):

```sh
$ curl -sf https://gobinaries.com/tj/robo | sh
```

From source:

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
 given will be passed.

 For example suppose you have the following task:

```yml
aws:
  exec: ssh tools aws
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

 Commands are executed via `sh -c`, thus you may use shell features and
 positional variables, for example:

```yml
hello:
  command: echo "Hello ${1:-there}"
```

 Yields:

```
$ robo hello
Hello there

$ robo hello Tobi
Hello there Tobi
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

 Any arguments given are simply appended.

### Scripts

 Shell scripts may be used instead of inline commands:

```yml
hello:
  summary: some task
  script: path/to/script.sh
```

 If the script is executable, it is invoked directly, this allows you
 to use `#!`:

```
$ echo -e '#!/usr/bin/env ruby\nputs "yo"' > yo.rb
$ chmod +x yo.rb
$ cat > robo.yml
yo:
  summary: yo from rb
  script: yo.rb
^C
$ robo yo
yo
```

Script paths are relative to the _config_ file, not the working directory.

### Usage

 Tasks may optionally specify usage parameters, which display
 upon help output:

```yml
events:
  summary: send data to the "events" topic
  exec: docker run -it events
  usage: "[project-id] [rate]"
```

### Examples

 Tasks may optionally specify any number of example commands, which
 display upon help output:

```yml
events:
  summary: send data to the "events" topic
  exec: docker run -it events
  usage: "[project-id] [rate]"
  examples:
    - description: Send 25 events a second to gy2d
      command: robo events gy2d 25
```

### Variables

 Robo supports variables via the [text/template](http://golang.org/pkg/text/template/) package. All you have to do is define a map of `variables` and use `{{` `}}` to refer to them.

 Here's an example:

```yml
stage:
  summary: Run commands against stage.
  exec: ssh {{.hosts.stage}} -t robo

prod:
  summary: Run commands against prod.
  exec: ssh {{.hosts.prod}} -t robo

variables:
  hosts:
    prod: bastion-prod
    stage: bastion-stage
```

The variables section does also interpolate itself with its own data via `{{ .var }}` and allows shell like command 
expressions via `$(echo true)` to be executed first providing the output result as a variable. Note that variables are 
interpolated first and then command expressions are evaluated. This will allow you to reduce repetitive variable definitions and declarations. 

````bash
hash:
  summary: echos the git {{ .branch }} branch's git hash
  command: echo {{ .branch }} {{ .githash }}

variables:
  branch: master
  githash: $(git rev-parse --short {{ .branch }})
````

  Along with your own custom variables, robo defines the following variables:

```bash
$ robo variables

    robo.file: /Users/amir/dev/src/github.com/tj/robo/robo.yml
    robo.path: /Users/amir/dev/src/github.com/tj/robo

    user.home: /Users/amir
    user.name: Amir Abushareb
    user.username: amir

```

### Environment

Tasks may define `env` key with an array of environment variables, this allows you
to re-use robo configuration files, for example:

```yaml
// aws.yml
dev:
  summary: AWS commands in dev environment
  exec: aws
  env: ["AWS_PROFILE=eng-dev"]

stage:
  summary: AWS commands in stage environment
  exec: aws
  env: ["AWS_PROFILE=eng-stage"]

prod:
  summary: AWS commands in prod environment
  exec: aws
  env: ["AWS_PROFILE=eng-prod"]
```

You can also override environment variables:

```bash
$ cat > robo.yml
home:
  summary: overrides $HOME
  exec: echo $HOME
  env: ["HOME=/tmp"]
^C
$ robo home // => /tmp
```

Variables can also be used to set env vars.

```bash
$ cat > robo.yml
aws-stage:
  summary: AWS stage
  exec: aws
  env: ["AWS_PROFILE={{.aws.profile}}"]
variables:
  aws:
    profile: eng-stage
^C
$ robo aws-stage ...
```

Note that you cannot use shell featurs in the environment key.

### Setup / Cleanup
Some tasks or even your entire robo configuration may require steps upfront for setup or afterwards for a cleanup. The keywords `before` and `after` can be embedded into a task or into the overall robo configuration. It has the same executable syntax as a task: `script`, `exec` and `command`.
Defining it on a task level causes the steps to be executed before (respectively after) the task. Global before or after steps are invoked for _every_ task in the configuration.
All steps get interpolated the same way tasks and variables are interpolated.

```yaml
before:
  - command: echo "global before {{ .foo }}"
after:
  - script: /global/after-script.sh

foo:
  before:
    - command: echo "local before {{ .foo }}"
    - exec: git pull -r
  after:
    - command: echo "local after"
    - exec: git reset --hard HEAD
  exec: git status

variables:
  foo: bar
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
  exec: ssh stage-tools -t robo --config stage.yml
```

  You can also use robo's builtin variables `robo.path`, for example
  if you put all robofiles in together:

```bash
├── dev.yml
├── prod.yml
├── root.yml
└── stage.yml
```

  And you would like to call `dev`, `prod` and `stage` from `root`:

```yml
dev:
  summary: Development commands
  exec: robo --config {{ .robo.path }}/dev.yml

stage:
  ...
```

## Composition

  You can compose multiple commands into a single command
  by utilizing robo's built-in `robo.file` variable:

```yml
one:
  summary: echo one
  command: echo one

two:
  summary: echo two
  command: echo two

all:
  summary: echo one two
  command: |
    robo -c {{ .robo.file }} one
    robo -c {{ .robo.file }} two
```

```
$ robo all
one
two
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

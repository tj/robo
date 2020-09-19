
v0.7.0 / 2020-09-19
===================

  * robo(1): add --quiet option
  * cli: add ListNames()
  * make: adjust dist target to exclude darwin/386

v0.6.0 / 2020-09-03
===================

  * Introduce variable map self-interpolation and command variable interpolation. (@jenpet)

v0.5.5 / 2020-02-16
===================

  * task: fix script path resolve()
  * docs: add gobinaries.com install option

v0.5.4 / 2019-08-31
===================

  * examples: add user example
  * doc: add robo built-in variables doc
  * config: add user built-in variables

v0.5.3 / 2019-08-31
===================

  * task: override env vars on exec

v0.5.1 / 2019-08-31
===================

  * main: bump version
  * make: add dist target

v0.5.0 / 2019-08-31
===================

  * doc: add compose, chain, env
  * examples: add compose, chain, env, exec, executable scripts
  * cli: improve variables output
  * config: add robo vars
  * task: fix exec escape
  * env: interpolate env vars using variables
  * cmd: resolve --config path
  * add env var support to allow re-use of robo conf files
  * task: run script directly if it is executable
  * ci: bump to latest go version
  * go: migrate to go modules

v0.4.1 / 2016-03-29
===================

  * fix script path resolution

v0.4.0 / 2015-10-13
===================

  * fix panic on yaml parse error. Closes #10

v0.3.0 / 2015-05-07
===================

  * change command behaviour to allow for positional vars

v0.2.0 / 2015-04-12
===================

  * fix script path resolution, now relative to config file
  * add `variables` command to list defined variables
  * add vars to feature list
  * add History.md

v0.1.0 / 2015-04-12
===================

  * add variable support.

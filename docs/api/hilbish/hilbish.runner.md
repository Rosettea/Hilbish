---
title: Interface hilbish.runner
description: interactive command runner customization
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The runner interface contains functions that allow the user to change
how Hilbish interprets interactive input.
Users can add and change the default runner for interactive input to any
language or script of their choosing. A good example is using it to
write command in Fennel.

## Functions
### setMode(cb)
This is the same as the `hilbish.runnerMode` function. It takes a callback,
which will be used to execute all interactive input.
In normal cases, neither callbacks should be overrided by the user,
as the higher level functions listed below this will handle it.

### lua(cmd)
Evaluates `cmd` as Lua input. This is the same as using `dofile`
or `load`, but is appropriated for the runner interface.

### sh(cmd)
Runs a command in Hilbish's shell script interpreter.
This is the equivalent of using `source`.

## setCurrent(name)
Sets the current interactive/command line runner mode.

## add(name, runner)
Adds a runner to the table of available runners. If runner is a table,
it must have the run function in it.

## get(name)
Get a runner by name.

## set(name, runner)
Sets a runner by name. The runner table must have the run function in it.

## exec(cmd, runnerName)
Executes cmd with a runner. If runnerName isn't passed, it uses
the user's current runner.


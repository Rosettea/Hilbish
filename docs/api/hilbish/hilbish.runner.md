---
name: Interface hilbish.runner
description: interactive command runner customization
layout: apidoc
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


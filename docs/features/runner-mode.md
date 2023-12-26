---
title: Runner Mode
description: Customize the interactive script/command runner.
layout: doc
menu: 
  docs:
    parent: "Features"
---

Hilbish allows you to change how interactive text can be interpreted.
This is mainly due to the fact that the default method Hilbish uses
is that it runs Lua first and then falls back to shell script.

In some cases, someone might want to switch to just shell script to avoid
it while interactive but still have a Lua config, or go full Lua to use
Hilbish as a REPL. This also allows users to add alternative languages like
Fennel as the interactive script runner.

Runner mode can also be used to handle specific kinds of input before
evaluating like normal, which is how [Link.hsh](https://github.com/TorchedSammy/Link.hsh)
handles links.

The "runner mode" of Hilbish is customizable via `hilbish.runnerMode`,
which determines how Hilbish will run user input. By default, this is
set to `hybrid` which is the previously mentioned behaviour of running Lua
first then going to shell script. If you want the reverse order, you can
set it to `hybridRev` and for isolated modes there is `sh` and `lua`
respectively.

You can also set it to a function, which will be called everytime Hilbish
needs to run interactive input. For more detail, see the [API documentation](../../api/hilbish/hilbish.runner)

The `hilbish.runner` interface is an alternative to using `hilbish.runnerMode`
and also provides the shell script and Lua runner functions that Hilbish itself uses.

A runner function is expected to return a table with the following values:
- `exitCode` (number): Exit code of the command
- `input` (string): The text input of the user. This is used by Hilbish to append extra input, in case
more is requested.
- `err` (string): A string that represents an error from the runner.
This should only be set when, for example, there is a syntax error.
It can be set to a few special values for Hilbish to throw the right
hooks and have a better looking message.
	- `<command>: not-found` will throw a `command.not-found` hook
	based on what `<command>` is.
	- `<command>: not-executable` will throw a `command.not-executable` hook.
- `continue` (boolean): Whether Hilbish should prompt the user for no input

## Functions
These are the "low level" functions for the `hilbish.runner` interface.

+ setMode(mode) > The same as `hilbish.runnerMode`
+ sh(input) -> table > Runs `input` in Hilbish's sh interpreter
+ lua(input) -> table > Evals `input` as Lua code

These functions should be preferred over the previous ones.
+ setCurrent(mode) > The same as `setMode`, but works with runners managed
via the functions below.
+ add(name, runner) > Adds a runner to a table of available runners. The `runner`
argument is either a function or a table with a run callback.
+ set(name, runner) > The same as `add` but requires passing a table and
overwrites if the `name`d runner already exists.
+ get(name) > runner > Gets a runner by name. It is a table with at least a
run function, to run input.
+ exec(cmd, runnerName) > Runs `cmd` with a runner. If `runnerName` isn't passed,
the current runner mode is used.

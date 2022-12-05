---
name: Module hilbish
description: the core Hilbish API
layout: apidoc
---

## Introduction
The Hilbish module includes the core API, containing
interfaces and functions which directly relate to shell functionality.

## Functions
### alias(cmd, orig)
Sets an alias of `cmd` to `orig`

### appendPath(dir)
Appends `dir` to $PATH

### complete(scope, cb)
Registers a completion handler for `scope`.
A `scope` is currently only expected to be `command.<cmd>`,
replacing <cmd> with the name of the command (for example `command.git`).
`cb` must be a function that returns a table of "completion groups."
Check `doc completions` for more information.

### cwd()
Returns the current directory of the shell

### exec(cmd)
Replaces running hilbish with `cmd`

### goro(fn)
Puts `fn` in a goroutine

### highlighter(line)
Line highlighter handler. This is mainly for syntax highlighting, but in
reality could set the input of the prompt to *display* anything. The
callback is passed the current line and is expected to return a line that
will be used as the input display.

### hinter(line, pos)
The command line hint handler. It gets called on every key insert to
determine what text to use as an inline hint. It is passed the current
line and cursor position. It is expected to return a string which is used
as the text for the hint. This is by default a shim. To set hints,
override this function with your custom handler.

### inputMode(mode)
Sets the input mode for Hilbish's line reader. Accepts either emacs or vim

### interval(cb, time)
Runs the `cb` function every `time` milliseconds.
Returns a `timer` object (see `doc timers`).

### multiprompt(str)
Changes the continued line prompt to `str`

### prependPath(dir)
Prepends `dir` to $PATH

### prompt(str, typ)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device

### read(prompt) -> input
Read input from the user, using Hilbish's line editor/input reader.
This is a separate instance from the one Hilbish actually uses.
Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)

### run(cmd, returnOut) -> exitCode, stdout, stderr
Runs `cmd` in Hilbish's sh interpreter.
If returnOut is true, the outputs of `cmd` will be returned as the 2nd and
3rd values instead of being outputted to the terminal.

### runnerMode(mode)
Sets the execution/runner mode for interactive Hilbish. This determines whether
Hilbish wll try to run input as Lua and/or sh or only do one of either.
Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
sh, and lua. It also accepts a function, to which if it is passed one
will call it to execute user input instead.

### timeout(cb, time)
Runs the `cb` function after `time` in milliseconds
Returns a `timer` object (see `doc timers`).

### which(name)
Checks if `name` is a valid command


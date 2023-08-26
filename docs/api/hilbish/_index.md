---
title: Module hilbish
description: the core Hilbish API
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The Hilbish module includes the core API, containing
interfaces and functions which directly relate to shell functionality.

## Interface fields
- `ver`: The version of Hilbish
- `goVersion`: The version of Go that Hilbish was compiled with
- `user`: Username of the user
- `host`: Hostname of the machine
- `dataDir`: Directory for Hilbish data files, including the docs and default modules
- `interactive`: Is Hilbish in an interactive shell?
- `login`: Is Hilbish the login shell?
- `vimMode`: Current Vim input mode of Hilbish (will be nil if not in Vim input mode)
- `exitCode`: xit code of the last executed command

## Functions
### hilbish.alias(cmd, orig)
Sets an alias of `cmd` to `orig`
#### Parameters
This function has no parameters.  

### hilbish.appendPath(dir)
Appends `dir` to $PATH
#### Parameters
This function has no parameters.  

### hilbish.complete(scope, cb)
Registers a completion handler for `scope`.
A `scope` is currently only expected to be `command.<cmd>`,
replacing <cmd> with the name of the command (for example `command.git`).
`cb` must be a function that returns a table of "completion groups."
Check `doc completions` for more information.
#### Parameters
This function has no parameters.  

### hilbish.cwd() -> string
Returns the current directory of the shell
#### Parameters
This function has no parameters.  

### hilbish.exec(cmd)
Replaces running hilbish with `cmd`
#### Parameters
This function has no parameters.  

### hilbish.goro(fn)
Puts `fn` in a goroutine
#### Parameters
This function has no parameters.  

### hilbish.highlighter(line)
Line highlighter handler. This is mainly for syntax highlighting, but in
reality could set the input of the prompt to *display* anything. The
callback is passed the current line and is expected to return a line that
will be used as the input display.
Note that to set a highlighter, one has to override this function.
Example:
```
function hilbish.highlighter(line)
   return line:gsub('"%w+"', function(c) return lunacolors.green(c) end)
end
```
This code will highlight all double quoted strings in green.
#### Parameters
This function has no parameters.  

### hilbish.hinter(line, pos)
The command line hint handler. It gets called on every key insert to
determine what text to use as an inline hint. It is passed the current
line and cursor position. It is expected to return a string which is used
as the text for the hint. This is by default a shim. To set hints,
override this function with your custom handler.
#### Parameters
This function has no parameters.  

### hilbish.inputMode(mode)
Sets the input mode for Hilbish's line reader. Accepts either emacs or vim
#### Parameters
This function has no parameters.  

### hilbish.interval(cb, time) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
Runs the `cb` function every `time` milliseconds.
This creates a timer that starts immediately.
#### Parameters
This function has no parameters.  

### hilbish.multiprompt(str)
Changes the continued line prompt to `str`
#### Parameters
This function has no parameters.  

### hilbish.prependPath(dir)
Prepends `dir` to $PATH
#### Parameters
This function has no parameters.  

### hilbish.prompt(str, typ)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device
#### Parameters
This function has no parameters.  

### hilbish.read(prompt) -> input (string)
Read input from the user, using Hilbish's line editor/input reader.
This is a separate instance from the one Hilbish actually uses.
Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
#### Parameters
This function has no parameters.  

### hilbish.run(cmd, returnOut) -> exitCode (number), stdout (string), stderr (string)
Runs `cmd` in Hilbish's sh interpreter.
If returnOut is true, the outputs of `cmd` will be returned as the 2nd and
3rd values instead of being outputted to the terminal.
#### Parameters
This function has no parameters.  

### hilbish.runnerMode(mode)
Sets the execution/runner mode for interactive Hilbish. This determines whether
Hilbish wll try to run input as Lua and/or sh or only do one of either.
Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
sh, and lua. It also accepts a function, to which if it is passed one
will call it to execute user input instead.
#### Parameters
This function has no parameters.  

### hilbish.timeout(cb, time) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
Runs the `cb` function after `time` in milliseconds.
This creates a timer that starts immediately.
#### Parameters
This function has no parameters.  

### hilbish.which(name) -> string
Checks if `name` is a valid command.
Will return the path of the binary, or a basename if it's a commander.
#### Parameters
This function has no parameters.  

## Types
## Sink
A sink is a structure that has input and/or output to/from
a desination.

### Methods
#### autoFlush(auto)
Sets/toggles the option of automatically flushing output.
A call with no argument will toggle the value.

#### flush()
Flush writes all buffered input to the sink.

#### read() -> string
Reads input from the sink.

#### write(str)
Writes data to a sink.

#### writeln(str)
Writes data to a sink with a newline at the end.


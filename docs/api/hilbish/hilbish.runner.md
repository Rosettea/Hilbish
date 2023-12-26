---
title: Module hilbish.runner
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

Runners are functions that evaluate user input. The default runners in
Hilbish can run shell script and Lua code.

A runner is passed the input and has to return a table with these values.
All are not required, only the useful ones the runner needs to return.
(So if there isn't an error, just omit `err`.)

- `exitCode` (number): A numerical code to indicate the exit result.
- `input` (string): The user input. This will be used to add
to the history.
- `err` (string): A string to indicate an interal error for the runner.
It can be set to a few special values for Hilbish to throw the right hooks and have a better looking message:

`[command]: not-found` will throw a command.not-found hook based on what `[command]` is.

`[command]: not-executable` will throw a command.not-executable hook.
- `continue` (boolean): Whether to prompt the user for more input.

Here is a simple example of a fennel runner. It falls back to
shell script if fennel eval has an error.
```lua
local fennel = require 'fennel'

hilbish.runnerMode(function(input)
	local ok = pcall(fennel.eval, input)
	if ok then
		return {
			input = input
		}
	end

	return hilbish.runner.sh(input)
end)
```

## Functions
|||
|----|----|
|<a href="#runner.setMode">setMode(cb)</a>|This is the same as the `hilbish.runnerMode` function.|
|<a href="#runner.lua">lua(cmd)</a>|Evaluates `cmd` as Lua input. This is the same as using `dofile`|
|<a href="#runner.sh">sh(cmd)</a>|Runs a command in Hilbish's shell script interpreter.|

<hr>
<div id='runner.setMode'>
<h4 class='heading'>
hilbish.runner.setMode(cb)
<a href="#runner.setMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

This is the same as the `hilbish.runnerMode` function.  
It takes a callback, which will be used to execute all interactive input.  
In normal cases, neither callbacks should be overrided by the user,  
as the higher level functions listed below this will handle it.  

#### Parameters
`function` **`cb`**  


</div>

<hr>
<div id='runner.lua'>
<h4 class='heading'>
hilbish.runner.lua(cmd)
<a href="#runner.lua" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Evaluates `cmd` as Lua input. This is the same as using `dofile`  
or `load`, but is appropriated for the runner interface.  

#### Parameters
`string` **`cmd`**  


</div>

<hr>
<div id='runner.sh'>
<h4 class='heading'>
hilbish.runner.sh(cmd)
<a href="#runner.sh" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs a command in Hilbish's shell script interpreter.  
This is the equivalent of using `source`.  

#### Parameters
`string` **`cmd`**  


</div>


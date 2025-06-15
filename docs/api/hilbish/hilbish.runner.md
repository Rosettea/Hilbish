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

- `exitCode` (number): Exit code of the command
- `input` (string): The text input of the user. This is used by Hilbish to append extra input, in case
more is requested.
- `err` (string): A string that represents an error from the runner.
This should only be set when, for example, there is a syntax error.
It can be set to a few special values for Hilbish to throw the right
hooks and have a better looking message.
	- `\<command>: not-found` will throw a `command.not-found` hook
	based on what `\<command>` is.
	- `\<command>: not-executable` will throw a `command.not-executable` hook.
- `continue` (boolean): Whether Hilbish should prompt the user for no input
- `newline` (boolean): Whether a newline should be added at the end of `input`.

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
|<a href="#runner.lua">lua(cmd)</a>|Evaluates `cmd` as Lua input. This is the same as using `dofile`|
|<a href="#sh">sh()</a>|nil|
|<a href="#setMode">setMode(mode)</a>|**NOTE: This function is deprecated and will be removed in 3.0**|
|<a href="#setCurrent">setCurrent(name)</a>|Sets Hilbish's runner mode by name.|
|<a href="#set">set(name, runner)</a>|*Sets* a runner by name. The difference between this function and|
|<a href="#run">run(input, priv)</a>|Runs `input` with the currently set Hilbish runner.|
|<a href="#getCurrent">getCurrent()</a>|Returns the current runner by name.|
|<a href="#get">get(name)</a>|Get a runner by name.|
|<a href="#exec">exec(cmd, runnerName)</a>|Executes `cmd` with a runner.|
|<a href="#add">add(name, runner)</a>|Adds a runner to the table of available runners.|

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
<div id='add'>
<h4 class='heading'>
hilbish.runner.add(name, runner)
<a href="#add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Adds a runner to the table of available runners.
If runner is a table, it must have the run function in it.
#### Parameters
`name` **`string`**  
 Name of the runner

`runner` **`function|table`**  
 

</div>

<hr>
<div id='exec'>
<h4 class='heading'>
hilbish.runner.exec(cmd, runnerName)
<a href="#exec" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Executes `cmd` with a runner.
If `runnerName` is not specified, it uses the default Hilbish runner.
#### Parameters
`cmd` **`string`**  


`runnerName` **`string?`**  


</div>

<hr>
<div id='get'>
<h4 class='heading'>
hilbish.runner.get(name)
<a href="#get" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Get a runner by name.
#### Parameters
`name` **`string`**  
 Name of the runner to retrieve.

</div>

<hr>
<div id='getCurrent'>
<h4 class='heading'>
hilbish.runner.getCurrent()
<a href="#getCurrent" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the current runner by name.
#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='run'>
<h4 class='heading'>
hilbish.runner.run(input, priv)
<a href="#run" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs `input` with the currently set Hilbish runner.
This method is how Hilbish executes commands.
`priv` is an optional boolean used to state if the input should be saved to history.
#### Parameters
`input` **`string`**  


`priv` **`bool`**  


</div>

<hr>
<div id='set'>
<h4 class='heading'>
hilbish.runner.set(name, runner)
<a href="#set" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

*Sets* a runner by name. The difference between this function and
add, is set will *not* check if the named runner exists.
The runner table must have the run function in it.
#### Parameters
`name` **`string`**  


`runner` **`table`**  


</div>

<hr>
<div id='setCurrent'>
<h4 class='heading'>
hilbish.runner.setCurrent(name)
<a href="#setCurrent" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets Hilbish's runner mode by name.
#### Parameters
`name` **`string`**  


</div>

<hr>
<div id='setMode'>
<h4 class='heading'>
hilbish.runner.setMode(mode)
<a href="#setMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

**NOTE: This function is deprecated and will be removed in 3.0**
Use `hilbish.runner.setCurrent` instead.
This is the same as the `hilbish.runnerMode` function.
It takes a callback, which will be used to execute all interactive input.
Or a string which names the runner mode to use.
#### Parameters
`mode` **`string|function`**  


</div>

<hr>
<div id='sh'>
<h4 class='heading'>
hilbish.runner.sh()
<a href="#sh" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>


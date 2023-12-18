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

## Functions
|||
|----|----|
|<a href="#alias">alias(cmd, orig)</a>|Sets an alias, with a name of `cmd` to another command.|
|<a href="#appendPath">appendPath(dir)</a>|Appends the provided dir to the command path (`$PATH`)|
|<a href="#complete">complete(scope, cb)</a>|Registers a completion handler for the specified scope.|
|<a href="#cwd">cwd() -> string</a>|Returns the current directory of the shell|
|<a href="#exec">exec(cmd)</a>|Replaces the currently running Hilbish instance with the supplied command.|
|<a href="#goro">goro(fn)</a>|Puts `fn` in a Goroutine.|
|<a href="#highlighter">highlighter(line)</a>|Line highlighter handler.|
|<a href="#hinter">hinter(line, pos)</a>|The command line hint handler. It gets called on every key insert to|
|<a href="#inputMode">inputMode(mode)</a>|Sets the input mode for Hilbish's line reader. Accepts either emacs or vim.|
|<a href="#interval">interval(cb, time) -> @Timer</a>|Runs the `cb` function every `time` milliseconds.|
|<a href="#multiprompt">multiprompt(str)</a>|Changes the text prompt when Hilbish asks for more input.|
|<a href="#prependPath">prependPath(dir)</a>|Prepends `dir` to $PATH.|
|<a href="#prompt">prompt(str, typ)</a>|Changes the shell prompt to the provided string.|
|<a href="#read">read(prompt) -> input (string)</a>|Read input from the user, using Hilbish's line editor/input reader.|
|<a href="#run">run(cmd, returnOut) -> exitCode (number), stdout (string), stderr (string)</a>|Runs `cmd` in Hilbish's shell script interpreter.|
|<a href="#runnerMode">runnerMode(mode)</a>|Sets the execution/runner mode for interactive Hilbish. This determines whether|
|<a href="#timeout">timeout(cb, time) -> @Timer</a>|Runs the `cb` function after `time` in milliseconds.|
|<a href="#which">which(name) -> string</a>|Checks if `name` is a valid command.|

## Static module fields
|||
|----|----|
|ver|The version of Hilbish|
|goVersion|The version of Go that Hilbish was compiled with|
|user|Username of the user|
|host|Hostname of the machine|
|dataDir|Directory for Hilbish data files, including the docs and default modules|
|interactive|Is Hilbish in an interactive shell?|
|login|Is Hilbish the login shell?|
|vimMode|Current Vim input mode of Hilbish (will be nil if not in Vim input mode)|
|exitCode|Exit code of the last executed command|

<hr><div id='alias'>
<h4 class='heading'>
hilbish.alias(cmd, orig)
<a href="#alias" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets an alias, with a name of `cmd` to another command.  
  
  
#### Parameters
`string` **`cmd`**  
Name of the alias

`string` **`orig`**  
Command that will be aliased

#### Example
```lua
-- With this, "ga file" will turn into "git add file"
hilbish.alias('ga', 'git add')

-- Numbered substitutions are supported here!
hilbish.alias('dircount', 'ls %1 | wc -l')
-- "dircount ~" would count how many files are in ~ (home directory).
```
</div>

<hr><div id='appendPath'>
<h4 class='heading'>
hilbish.appendPath(dir)
<a href="#appendPath" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Appends the provided dir to the command path (`$PATH`)  
  
  
#### Parameters
`string|table` **`dir`**  
Directory (or directories) to append to path

#### Example
```lua
hilbish.appendPath '~/go/bin'
-- Will add ~/go/bin to the command path.

-- Or do multiple:
hilbush.appendPath {
	'~/go/bin',
	'~/.local/bin'
}
```
</div>

<hr><div id='complete'>
<h4 class='heading'>
hilbish.complete(scope, cb)
<a href="#complete" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Registers a completion handler for the specified scope.  
A `scope` is currently only expected to be `command.<cmd>`,  
replacing <cmd> with the name of the command (for example `command.git`).  
The documentation for completions, under Features/Completions or `doc completions`  
provides more details.  
#### Parameters
`string` **`scope`**  


`function` **`cb`**  


</div>

<hr><div id='cwd'>
<h4 class='heading'>
hilbish.cwd() -> string
<a href="#cwd" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the current directory of the shell  
#### Parameters
This function has no parameters.  
</div>

<hr><div id='exec'>
<h4 class='heading'>
hilbish.exec(cmd)
<a href="#exec" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Replaces the currently running Hilbish instance with the supplied command.  
This can be used to do an in-place restart.  
#### Parameters
`string` **`cmd`**  


</div>

<hr><div id='goro'>
<h4 class='heading'>
hilbish.goro(fn)
<a href="#goro" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Puts `fn` in a Goroutine.  
This can be used to run any function in another thread.  
**NOTE: THIS FUNCTION MAY CRASH HILBISH IF OUTSIDE VARIABLES ARE ACCESSED.**  
#### Parameters
`function` **`fn`**  


</div>

<hr><div id='highlighter'>
<h4 class='heading'>
hilbish.highlighter(line)
<a href="#highlighter" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Line highlighter handler.  
This is mainly for syntax highlighting, but in reality could set the input  
of the prompt to *display* anything. The callback is passed the current line  
and is expected to return a line that will be used as the input display.  
Note that to set a highlighter, one has to override this function.  
  
#### Parameters
`string` **`line`**  


#### Example
```lua
--This code will highlight all double quoted strings in green.
function hilbish.highlighter(line)
   return line:gsub('"%w+"', function(c) return lunacolors.green(c) end)
end
```
</div>

<hr><div id='hinter'>
<h4 class='heading'>
hilbish.hinter(line, pos)
<a href="#hinter" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

The command line hint handler. It gets called on every key insert to  
determine what text to use as an inline hint. It is passed the current  
line and cursor position. It is expected to return a string which is used  
as the text for the hint. This is by default a shim. To set hints,  
override this function with your custom handler.  
  
  
#### Parameters
`string` **`line`**  


`number` **`pos`**  


#### Example
```lua
-- this will display "hi" after the cursor in a dimmed color.
function hilbish.hinter(line, pos)
	return 'hi'
end
```
</div>

<hr><div id='inputMode'>
<h4 class='heading'>
hilbish.inputMode(mode)
<a href="#inputMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets the input mode for Hilbish's line reader. Accepts either emacs or vim.  
`emacs` is the default. Setting it to `vim` changes behavior of input to be  
Vim-like with modes and Vim keybinds.  
#### Parameters
`string` **`mode`**  


</div>

<hr><div id='interval'>
<h4 class='heading'>
hilbish.interval(cb, time) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
<a href="#interval" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs the `cb` function every `time` milliseconds.  
This creates a timer that starts immediately.  
#### Parameters
`function` **`cb`**  


`number` **`time`**  


</div>

<hr><div id='multiprompt'>
<h4 class='heading'>
hilbish.multiprompt(str)
<a href="#multiprompt" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Changes the text prompt when Hilbish asks for more input.  
This will show up when text is incomplete, like a missing quote  
  
  
#### Parameters
`string` **`str`**  


#### Example
```lua
--[[
imagine this is your text input:
user ~ ∆ echo "hey

but there's a missing quote! hilbish will now prompt you so the terminal
will look like:
user ~ ∆ echo "hey
--> ...!"

so then you get
user ~ ∆ echo "hey
--> ...!"
hey ...!
]]--
hilbish.multiprompt '-->'
```
</div>

<hr><div id='prependPath'>
<h4 class='heading'>
hilbish.prependPath(dir)
<a href="#prependPath" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Prepends `dir` to $PATH.  
#### Parameters
`string` **`dir`**  


</div>

<hr><div id='prompt'>
<h4 class='heading'>
hilbish.prompt(str, typ)
<a href="#prompt" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Changes the shell prompt to the provided string.  
There are a few verbs that can be used in the prompt text.  
These will be formatted and replaced with the appropriate values.  
`%d` - Current working directory  
`%u` - Name of current user  
`%h` - Hostname of device  
  
#### Parameters
`string` **`str`**  


`string` **`typ?`**  
Type of prompt, being left or right. Left by default.

#### Example
```lua
-- the default hilbish prompt without color
hilbish.prompt '%u %d ∆'
-- or something of old:
hilbish.prompt '%u@%h :%d $'
-- prompt: user@hostname: ~/directory $
```
</div>

<hr><div id='read'>
<h4 class='heading'>
hilbish.read(prompt) -> input (string)
<a href="#read" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Read input from the user, using Hilbish's line editor/input reader.  
This is a separate instance from the one Hilbish actually uses.  
Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen).  
#### Parameters
`string` **`prompt?`**  


</div>

<hr><div id='run'>
<h4 class='heading'>
hilbish.run(cmd, returnOut) -> exitCode (number), stdout (string), stderr (string)
<a href="#run" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs `cmd` in Hilbish's shell script interpreter.  
#### Parameters
`string` **`cmd`**  


`boolean` **`returnOut`**  
If this is true, the function will return the standard output and error of the command instead of printing it.

</div>

<hr><div id='runnerMode'>
<h4 class='heading'>
hilbish.runnerMode(mode)
<a href="#runnerMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets the execution/runner mode for interactive Hilbish. This determines whether  
Hilbish wll try to run input as Lua and/or sh or only do one of either.  
Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),  
sh, and lua. It also accepts a function, to which if it is passed one  
will call it to execute user input instead.  
#### Parameters
`string|function` **`mode`**  


</div>

<hr><div id='timeout'>
<h4 class='heading'>
hilbish.timeout(cb, time) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
<a href="#timeout" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs the `cb` function after `time` in milliseconds.  
This creates a Timer that starts immediately.  
#### Parameters
`function` **`cb`**  


`number` **`time`**  


</div>

<hr><div id='which'>
<h4 class='heading'>
hilbish.which(name) -> string
<a href="#which" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Checks if `name` is a valid command.  
Will return the path of the binary, or a basename if it's a commander.  
#### Parameters
`string` **`name`**  


</div>

## Types
<hr>

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
Reads a liine of input from the sink.

#### readAll() -> string
Reads all input from the sink.

#### write(str)
Writes data to a sink.

#### writeln(str)
Writes data to a sink with a newline at the end.


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

## Functions
|||
|----|----|
|<a href="#runner.setMode">setMode(cb)</a>|This is the same as the `hilbish.runnerMode` function. It takes a callback,|
|<a href="#runner.lua">lua(cmd)</a>|Evaluates `cmd` as Lua input. This is the same as using `dofile`|
|<a href="#runner.sh">sh(cmd)</a>|Runs a command in Hilbish's shell script interpreter.|

<hr><div id='runner.setMode'>
<h4 class='heading'>
hilbish.runner.setMode(cb)
<a href="#runner.setMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

This is the same as the `hilbish.runnerMode` function. It takes a callback,
which will be used to execute all interactive input.
In normal cases, neither callbacks should be overrided by the user,
as the higher level functions listed below this will handle it.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='runner.lua'>
<h4 class='heading'>
hilbish.runner.lua(cmd)
<a href="#runner.lua" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Evaluates `cmd` as Lua input. This is the same as using `dofile`
or `load`, but is appropriated for the runner interface.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='runner.sh'>
<h4 class='heading'>
hilbish.runner.sh(cmd)
<a href="#runner.sh" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Runs a command in Hilbish's shell script interpreter.
This is the equivalent of using `source`.
#### Parameters
This function has no parameters.  
</div>


---
title: Module hilbish.history
description: command history
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The history interface deals with command history.
This includes the ability to override functions to change the main
method of saving history.

## Functions
|||
|----|----|
|<a href="#history.add">add(cmd)</a>|Adds a command to the history.|
|<a href="#history.all">all() -> table</a>|Retrieves all history.|
|<a href="#history.clear">clear()</a>|Deletes all commands from the history.|
|<a href="#history.get">get(idx)</a>|Retrieves a command from the history based on the `idx`.|
|<a href="#history.size">size() -> number</a>|Returns the amount of commands in the history.|

<hr><div id='history.add'>
<h4 class='heading'>
hilbish.history.add(cmd)
<a href="#history.add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Adds a command to the history.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='history.all'>
<h4 class='heading'>
hilbish.history.all() -> table
<a href="#history.all" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Retrieves all history.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='history.clear'>
<h4 class='heading'>
hilbish.history.clear()
<a href="#history.clear" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Deletes all commands from the history.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='history.get'>
<h4 class='heading'>
hilbish.history.get(idx)
<a href="#history.get" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Retrieves a command from the history based on the `idx`.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='history.size'>
<h4 class='heading'>
hilbish.history.size() -> number
<a href="#history.size" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the amount of commands in the history.
#### Parameters
This function has no parameters.  
</div>


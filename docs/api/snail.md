---
title: Module snail
description: shell script interpreter library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

The snail library houses Hilbish's Lua wrapper of its shell script interpreter.
It's not very useful other than running shell scripts, which can be done with other
Hilbish functions.

## Functions
|||
|----|----|
|<a href="#new">new() -> @Snail</a>|Creates a new Snail instance.|

<hr>
<div id='new'>
<h4 class='heading'>
snail.new() -> <a href="/Hilbish/docs/api/snail/#snail" style="text-decoration: none;" id="lol">Snail</a>
<a href="#new" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new Snail instance.  

#### Parameters
This function has no parameters.  
</div>

## Types
<hr>

## Snail
A Snail is a shell script interpreter instance.

### Methods
#### dir(path)
Changes the directory of the snail instance.
The interpreter keeps its set directory even when the Hilbish process changes
directory, so this should be called on the `hilbish.cd` hook.

#### run(command, streams)
Runs a shell command. Works the same as `hilbish.run`, but only accepts a table of streams.


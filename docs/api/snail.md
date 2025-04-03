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
|<a href="#handleStream"></a>||
|<a href="#loaderFunc"></a>||
|<a href="#snailUserData"></a>||
|<a href="#snew">new() -> @Snail</a>|Creates a new Snail instance.|
|<a href="#splitInput"></a>||
|<a href="#Run"></a>||

<hr>
<div id='handleStream'>
<h4 class='heading'>
snail.
<a href="#handleStream" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='loaderFunc'>
<h4 class='heading'>
snail.
<a href="#loaderFunc" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='snailUserData'>
<h4 class='heading'>
snail.
<a href="#snailUserData" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='snew'>
<h4 class='heading'>
snail.new() -> <a href="/Hilbish/docs/api/snail/#snail" style="text-decoration: none;" id="lol">Snail</a>
<a href="#snew" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new Snail instance.  

#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='splitInput'>
<h4 class='heading'>
snail.
<a href="#splitInput" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='Run'>
<h4 class='heading'>
snail.
<a href="#Run" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>


#### Parameters
This function has no parameters.  
</div>

## Types
<hr>

## Snail
A Snail is a shell script interpreter instance.

### Methods
#### run(command, streams)
Runs a shell command. Works the same as `hilbish.run`.


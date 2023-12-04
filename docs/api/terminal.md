---
title: Module terminal
description: low level terminal library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The terminal library is a simple and lower level library for certain terminal interactions.

## Functions
|||
|----|----|
|<a href="#restoreState">restoreState()</a>|Restores the last saved state of the terminal|
|<a href="#saveState">saveState()</a>|Saves the current state of the terminal.|
|<a href="#setRaw">setRaw()</a>|Puts the terminal into raw mode.|
|<a href="#size">size()</a>|Gets the dimensions of the terminal. Returns a table with `width` and `height`|

<hr><div id='restoreState'>
<h4 class='heading'>
terminal.restoreState()
<a href="#restoreState" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Restores the last saved state of the terminal  
#### Parameters
This function has no parameters.  
</div>

<hr><div id='saveState'>
<h4 class='heading'>
terminal.saveState()
<a href="#saveState" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Saves the current state of the terminal.  
#### Parameters
This function has no parameters.  
</div>

<hr><div id='setRaw'>
<h4 class='heading'>
terminal.setRaw()
<a href="#setRaw" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Puts the terminal into raw mode.  
#### Parameters
This function has no parameters.  
</div>

<hr><div id='size'>
<h4 class='heading'>
terminal.size()
<a href="#size" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Gets the dimensions of the terminal. Returns a table with `width` and `height`  
NOTE: The size refers to the amount of columns and rows of text that can fit in the terminal.  
#### Parameters
This function has no parameters.  
</div>


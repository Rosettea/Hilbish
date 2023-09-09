---
title: Module hilbish.editor
description: interactions for Hilbish's line reader
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The hilbish.editor interface provides functions to
directly interact with the line editor in use.

## Functions
|||
|----|----|
|<a href="#editor.getLine">getLine() -> string</a>|Returns the current input line.|
|<a href="#editor.getVimRegister">getVimRegister(register) -> string</a>|Returns the text that is at the register.|
|<a href="#editor.insert">insert(text)</a>|Inserts text into the line.|
|<a href="#editor.setVimRegister">setVimRegister(register, text)</a>|Sets the vim register at `register` to hold the passed text.|

<hr><div id='editor.getLine'>
<h4 class='heading'>
hilbish.editor.getLine() -> string
<a href="#editor.getLine" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the current input line.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='editor.getVimRegister'>
<h4 class='heading'>
hilbish.editor.getVimRegister(register) -> string
<a href="#editor.getVimRegister" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the text that is at the register.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='editor.insert'>
<h4 class='heading'>
hilbish.editor.insert(text)
<a href="#editor.insert" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Inserts text into the line.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='editor.setVimRegister'>
<h4 class='heading'>
hilbish.editor.setVimRegister(register, text)
<a href="#editor.setVimRegister" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets the vim register at `register` to hold the passed text.
#### Parameters
This function has no parameters.  
</div>


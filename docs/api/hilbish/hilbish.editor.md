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
|<a href="#editor.deleteByAmount">deleteByAmount(amount)</a>|Deletes characters in the line by the given amount.|
|<a href="#editor.getLine">getLine() -> string</a>|Returns the current input line.|
|<a href="#editor.getVimRegister">getVimRegister(register) -> string</a>|Returns the text that is at the register.|
|<a href="#editor.insert">insert(text)</a>|Inserts text into the Hilbish command line.|
|<a href="#editor.getChar">getChar() -> string</a>|Reads a keystroke from the user. This is in a format of something like Ctrl-L.|
|<a href="#editor.setVimRegister">setVimRegister(register, text)</a>|Sets the vim register at `register` to hold the passed text.|

<hr>
<div id='editor.deleteByAmount'>
<h4 class='heading'>
hilbish.editor.deleteByAmount(amount)
<a href="#editor.deleteByAmount" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Deletes characters in the line by the given amount.  

#### Parameters
`number` **`amount`**  


</div>

<hr>
<div id='editor.getLine'>
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

<hr>
<div id='editor.getVimRegister'>
<h4 class='heading'>
hilbish.editor.getVimRegister(register) -> string
<a href="#editor.getVimRegister" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the text that is at the register.  

#### Parameters
`string` **`register`**  


</div>

<hr>
<div id='editor.insert'>
<h4 class='heading'>
hilbish.editor.insert(text)
<a href="#editor.insert" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Inserts text into the Hilbish command line.  

#### Parameters
`string` **`text`**  


</div>

<hr>
<div id='editor.getChar'>
<h4 class='heading'>
hilbish.editor.getChar() -> string
<a href="#editor.getChar" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Reads a keystroke from the user. This is in a format of something like Ctrl-L.  

#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='editor.setVimRegister'>
<h4 class='heading'>
hilbish.editor.setVimRegister(register, text)
<a href="#editor.setVimRegister" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sets the vim register at `register` to hold the passed text.  

#### Parameters
`string` **`register`**  


`string` **`text`**  


</div>


---
title: Module readline
description: line reader library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The readline module is responsible for reading input from the user.
The readline module is what Hilbish uses to read input from the user,
including all the interactive features of Hilbish like history search,
syntax highlighting, everything. The global Hilbish readline instance
is usable at `hilbish.editor`.

## Functions
|||
|----|----|
|<a href="#New">new() -> @Readline</a>|Creates a new readline instance.|

<hr>
<div id='New'>
<h4 class='heading'>
readline.new() -> <a href="/Hilbish/docs/api/readline/#readline" style="text-decoration: none;" id="lol">Readline</a>
<a href="#New" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new readline instance.  

#### Parameters
This function has no parameters.  
</div>

## Types
<hr>

## Readline

### Methods
#### deleteByAmount(amount)
Deletes characters in the line by the given amount.

#### getLine() -> string
Returns the current input line.

#### getVimRegister(register) -> string
Returns the text that is at the register.

#### insert(text)
Inserts text into the Hilbish command line.

#### log(text)
Prints a message *before* the prompt without it being interrupted by user input.

#### read() -> string
Reads input from the user.

#### getChar() -> string
Reads a keystroke from the user. This is in a format of something like Ctrl-L.

#### setVimRegister(register, text)
Sets the vim register at `register` to hold the passed text.


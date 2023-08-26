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
### hilbish.editor.getLine() -> string
Returns the current input line.
#### Parameters
This function has no parameters.  

### hilbish.editor.getVimRegister(register) -> string
Returns the text that is at the register.
#### Parameters
This function has no parameters.  

### hilbish.editor.insert(text)
Inserts text into the line.
#### Parameters
This function has no parameters.  

### hilbish.editor.setVimRegister(register, text)
Sets the vim register at `register` to hold the passed text.
#### Parameters
This function has no parameters.  


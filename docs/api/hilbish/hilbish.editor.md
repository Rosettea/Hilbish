---
title: Interface hilbish.editor
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
### getLine() -> string
Returns the current input line.

### getVimRegister(register) -> string
Returns the text that is at the register.

### insert(text)
Inserts text into the line.

### getChar() -> string
Reads a keystroke from the user. This is in a format
of something like Ctrl-L.

### setVimRegister(register, text)
Sets the vim register at `register` to hold the passed text.


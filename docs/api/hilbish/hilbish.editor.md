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
### getLine()
Returns the current input line.

### getVimRegister(register)
Returns the text that is at the register.

### insert(text)
Inserts text into the line.

### setVimRegister(register, text)
Sets the vim register at `register` to hold the passed text.


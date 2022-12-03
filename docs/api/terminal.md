---
name: Module terminal
description: low level terminal library
layout: apidoc
---

## Introduction
The terminal library is a simple and lower level library for certain terminal interactions.

## Functions
### size()
Gets the dimensions of the terminal. Returns a table with `width` and `height`
Note: this is not the size in relation to the dimensions of the display

### restoreState()
Restores the last saved state of the terminal

### saveState()
Saves the current state of the terminal

### setRaw()
Puts the terminal in raw mode


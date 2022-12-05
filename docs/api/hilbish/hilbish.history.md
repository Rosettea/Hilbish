---
name: Interface hilbish.history
description: command history
layout: apidoc
---

## Introduction
The history interface deals with command history.
This includes the ability to override functions to change the main
method of saving history.

## Functions
### add(cmd)
Adds a command to the history.

### clear()
Deletes all commands from the history.

### get(idx)
Retrieves a command from the history based on the `idx`.

### size()
Returns the amount of commands in the history.


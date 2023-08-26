---
title: Module hilbish.completions
description: tab completions
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The completions interface deals with tab completions.

## Functions
### hilbish.completions.call(name, query, ctx, fields) -> completionGroups (table), prefix (string)
Calls a completer function. This is mainly used to call
a command completer, which will have a `name` in the form
of `command.name`, example: `command.git`.
You can check `doc completions` for info on the `completionGroups` return value.
#### Parameters
This function has no parameters.  

### hilbish.completions.handler(line, pos)
The handler function is the callback for tab completion in Hilbish.
You can check the completions doc for more info.
#### Parameters
This function has no parameters.  

### hilbish.completions.bins(query, ctx, fields) -> entries (table), prefix (string)
Returns binary/executale completion candidates based on the provided query.
#### Parameters
This function has no parameters.  

### hilbish.completions.files(query, ctx, fields) -> entries (table), prefix (string)
Returns file completion candidates based on the provided query.
#### Parameters
This function has no parameters.  


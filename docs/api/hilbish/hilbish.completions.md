---
title: Interface hilbish.completions
description: tab completions
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The completions interface deals with tab completions.

## Functions
### call(name, query, ctx, fields)
Calls a completer function. This is mainly used to call
a command completer, which will have a `name` in the form
of `command.name`, example: `command.git`

### handler(line, pos)
The handler function is the callback for tab completion in Hilbish.
You can check the completions doc for more info.

### bins(query, ctx, fields)
Returns binary/executale completion candidates based on the provided query.

### files(query, ctx, fields)
Returns file completion candidates based on the provided query.


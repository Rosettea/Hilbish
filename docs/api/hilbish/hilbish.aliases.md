---
title: Module hilbish.aliases
description: command aliasing
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The alias interface deals with all command aliases in Hilbish.

## Functions
### hilbish.aliases.add(alias, cmd)
This is an alias (ha) for the `hilbish.alias` function.
#### Parameters
This function has no parameters.  

### hilbish.aliases.delete(name)
Removes an alias.
#### Parameters
This function has no parameters.  

### hilbish.aliases.list() -> table\<string, string>
Get a table of all aliases, with string keys as the alias and the value as the command.
#### Parameters
This function has no parameters.  

### hilbish.aliases.resolve(alias) -> command (string)
Tries to resolve an alias to its command.
#### Parameters
This function has no parameters.  


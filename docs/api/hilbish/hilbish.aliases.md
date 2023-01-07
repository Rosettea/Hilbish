---
title: Interface hilbish.aliases
description: command aliasing
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The alias interface deals with all command aliases in Hilbish.

## Functions
### add(alias, cmd)
This is an alias (ha) for the `hilbish.alias` function.

### delete(name)
Removes an alias.

### list() -> table<string, string>
Get a table of all aliases, with string keys as the alias and the value as the command.

### resolve(alias) -> command (string)
Tries to resolve an alias to its command.


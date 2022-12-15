---
title: Module commander
description: library for custom commands
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
Commander is a library for writing custom commands in Lua.

## Functions
### deregister(name)
Deregisters any command registered with `name`

### register(name, cb)
Register a command with `name` that runs `cb` when ran


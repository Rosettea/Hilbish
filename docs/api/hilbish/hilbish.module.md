---
title: Interface hilbish.module
description: native module loading
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
 The hilbish.module interface provides a function
to load Hilbish plugins/modules.
Hilbish modules are Go-written plugins (see https://pkg.go.dev/plugin)
that are used to add functionality to Hilbish that cannot be written
in Lua for any reason.

To make a valid native module, the Go plugin
has to export a Loader function with a signature like so:
`func(*rt.Runtime) rt.Value`.

`rt` in this case refers to the Runtime type at
https://pkg.go.dev/github.com/arnodel/golua@master/runtime#Runtime

Hilbish uses this package as its Lua runtime. You will need to read
it to use it for a native plugin.

## Functions
### load(path)
Loads a module at the designated `path`.
It will throw if any error occurs.


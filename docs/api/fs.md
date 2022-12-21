---
title: Module fs
description: filesystem interaction and functionality library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The fs module provides easy and simple access to filesystem functions
and other things, and acts an addition to the Lua standard library's
I/O and filesystem functions.

## Functions
### abs(path)
Gives an absolute version of `path`.

### basename(path)
Gives the basename of `path`. For the rules,
see Go's filepath.Base

### cd(dir)
Changes directory to `dir`

### dir(path)
Returns the directory part of `path`. For the rules, see Go's
filepath.Dir

### glob(pattern)
Glob all files and directories that match the pattern.
For the rules, see Go's filepath.Glob

### join(...)
Takes paths and joins them together with the OS's
directory separator (forward or backward slash).

### mkdir(name, recursive)
Makes a directory called `name`. If `recursive` is true, it will create its parent directories.

### readdir(dir) -> {}
Returns a table of files in `dir`.

### stat(path) -> {}
Returns a table of info about the `path`.
It contains the following keys:
name (string) - Name of the path
size (number) - Size of the path
mode (string) - Permission mode in an octal format string (with leading 0)
isDir (boolean) - If the path is a directory


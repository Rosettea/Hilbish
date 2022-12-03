---
name: Module fs
description: filesystem interaction and functionality library
layout: apidoc
---

## Introduction
The fs module provides easy and simple access to filesystem functions
and other things, and acts an addition to the Lua standard library's
I/O and filesystem functions.

## Functions
### basename(path)
Gives the basename of `path`. For the rules,
see Go's filepath.Base

### cd(dir)
Changes directory to `dir`

### dir(path)
Returns the directory part of `path`. For the rules, see Go's
filepath.Dir

### readdir(dir)
Returns a table of files in `dir`

### abs(path)
Gives an absolute version of `path`.

### glob(pattern)
Glob all files and directories that match the pattern.
For the rules, see Go's filepath.Glob

### join(paths...)
Takes paths and joins them together with the OS's
directory separator (forward or backward slash).

### mkdir(name, recursive)
Makes a directory called `name`. If `recursive` is true, it will create its parent directories.

### stat(path)
Returns info about `path`


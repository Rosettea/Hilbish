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
### fs.abs(path) -> string
Gives an absolute version of `path`.
#### Parameters
This function has no parameters.  

### fs.basename(path) -> string
Gives the basename of `path`. For the rules,
see Go's filepath.Base
#### Parameters
This function has no parameters.  

### fs.cd(dir)
Changes directory to `dir`
#### Parameters
This function has no parameters.  

### fs.dir(path) -> string
Returns the directory part of `path`. For the rules, see Go's
filepath.Dir
#### Parameters
This function has no parameters.  

### fs.glob(pattern) -> matches (table)
Glob all files and directories that match the pattern.
For the rules, see Go's filepath.Glob
#### Parameters
This function has no parameters.  

### fs.join(...) -> string
Takes paths and joins them together with the OS's
directory separator (forward or backward slash).
#### Parameters
This function has no parameters.  

### fs.mkdir(name, recursive)
Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
#### Parameters
This function has no parameters.  

### fs.readdir(dir) -> {}
Returns a table of files in `dir`.
#### Parameters
This function has no parameters.  

### fs.stat(path) -> {}
Returns a table of info about the `path`.
It contains the following keys:
name (string) - Name of the path
size (number) - Size of the path
mode (string) - Permission mode in an octal format string (with leading 0)
isDir (boolean) - If the path is a directory
#### Parameters
This function has no parameters.  


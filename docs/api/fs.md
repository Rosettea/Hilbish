---
name: Module fs
description: 
layout: apidoc
---

## Introduction


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

### join(paths...)
Takes paths and joins them together with the OS's
directory separator (forward or backward slash).

### mkdir(name, recursive)
Makes a directory called `name`. If `recursive` is true, it will create its parent directories.

### readdir(dir)
Returns a table of files in `dir`

### stat(path)
Returns info about `path`


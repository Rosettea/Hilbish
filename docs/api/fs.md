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
|||
|----|----|
|<a href="#abs">abs(path) -> string</a>|Gives an absolute version of `path`.|
|<a href="#basename">basename(path) -> string</a>|Gives the basename of `path`. For the rules,|
|<a href="#cd">cd(dir)</a>|Changes directory to `dir`|
|<a href="#dir">dir(path) -> string</a>|Returns the directory part of `path`. For the rules, see Go's|
|<a href="#glob">glob(pattern) -> matches (table)</a>|Glob all files and directories that match the pattern.|
|<a href="#join">join(...) -> string</a>|Takes paths and joins them together with the OS's|
|<a href="#mkdir">mkdir(name, recursive)</a>|Makes a directory called `name`. If `recursive` is true, it will create its parent directories.|
|<a href="#readdir">readdir(dir) -> {}</a>|Returns a table of files in `dir`.|
|<a href="#stat">stat(path) -> {}</a>|Returns a table of info about the `path`.|

## Functions
<hr><div id='abs'>
<h4 class='heading'>
fs.abs(path) -> string
<a href="#abs" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Gives an absolute version of `path`.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='basename'>
<h4 class='heading'>
fs.basename(path) -> string
<a href="#basename" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Gives the basename of `path`. For the rules,
see Go's filepath.Base
#### Parameters
This function has no parameters.  
</div>

<hr><div id='cd'>
<h4 class='heading'>
fs.cd(dir)
<a href="#cd" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Changes directory to `dir`
#### Parameters
This function has no parameters.  
</div>

<hr><div id='dir'>
<h4 class='heading'>
fs.dir(path) -> string
<a href="#dir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the directory part of `path`. For the rules, see Go's
filepath.Dir
#### Parameters
This function has no parameters.  
</div>

<hr><div id='glob'>
<h4 class='heading'>
fs.glob(pattern) -> matches (table)
<a href="#glob" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Glob all files and directories that match the pattern.
For the rules, see Go's filepath.Glob
#### Parameters
This function has no parameters.  
</div>

<hr><div id='join'>
<h4 class='heading'>
fs.join(...) -> string
<a href="#join" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Takes paths and joins them together with the OS's
directory separator (forward or backward slash).
#### Parameters
This function has no parameters.  
</div>

<hr><div id='mkdir'>
<h4 class='heading'>
fs.mkdir(name, recursive)
<a href="#mkdir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='readdir'>
<h4 class='heading'>
fs.readdir(dir) -> {}
<a href="#readdir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns a table of files in `dir`.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='stat'>
<h4 class='heading'>
fs.stat(path) -> {}
<a href="#stat" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns a table of info about the `path`.
It contains the following keys:
name (string) - Name of the path
size (number) - Size of the path
mode (string) - Permission mode in an octal format string (with leading 0)
isDir (boolean) - If the path is a directory
#### Parameters
This function has no parameters.  
</div>


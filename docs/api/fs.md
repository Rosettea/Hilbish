---
title: Module fs
description: filesystem interaction and functionality library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

The fs module provides filesystem functions to Hilbish. While Lua's standard
library has some I/O functions, they're missing a lot of the basics. The `fs`
library offers more functions and will work on any operating system Hilbish does.

## Functions
|||
|----|----|
|<a href="#abs">abs(path) -> string</a>|Returns an absolute version of the `path`.|
|<a href="#basename">basename(path) -> string</a>|Returns the "basename," or the last part of the provided `path`. If path is empty,|
|<a href="#cd">cd(dir)</a>|Changes Hilbish's directory to `dir`.|
|<a href="#dir">dir(path) -> string</a>|Returns the directory part of `path`. If a file path like|
|<a href="#glob">glob(pattern) -> matches (table)</a>|Match all files based on the provided `pattern`.|
|<a href="#join">join(...path) -> string</a>|Takes any list of paths and joins them based on the operating system's path separator.|
|<a href="#mkdir">mkdir(name, recursive)</a>|Creates a new directory with the provided `name`.|
|<a href="#readdir">readdir(path) -> table[string]</a>|Returns a list of all files and directories in the provided path.|
|<a href="#stat">stat(path) -> {}</a>|Returns the information about a given `path`.|

## Static module fields
|||
|----|----|
|pathSep|The operating system's path separator.|

<hr>
<div id='abs'>
<h4 class='heading'>
fs.abs(path) -> string
<a href="#abs" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns an absolute version of the `path`.  
This can be used to resolve short paths like `..` to `/home/user`.  

#### Parameters
`string` **`path`**  


</div>

<hr>
<div id='basename'>
<h4 class='heading'>
fs.basename(path) -> string
<a href="#basename" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the "basename," or the last part of the provided `path`. If path is empty,  
`.` will be returned.  

#### Parameters
`string` **`path`**  
Path to get the base name of.

</div>

<hr>
<div id='cd'>
<h4 class='heading'>
fs.cd(dir)
<a href="#cd" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Changes Hilbish's directory to `dir`.  

#### Parameters
`string` **`dir`**  
Path to change directory to.

</div>

<hr>
<div id='dir'>
<h4 class='heading'>
fs.dir(path) -> string
<a href="#dir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the directory part of `path`. If a file path like  
`~/Documents/doc.txt` then this function will return `~/Documents`.  

#### Parameters
`string` **`path`**  
Path to get the directory for.

</div>

<hr>
<div id='glob'>
<h4 class='heading'>
fs.glob(pattern) -> matches (table)
<a href="#glob" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Match all files based on the provided `pattern`.  
For the syntax' refer to Go's filepath.Match function: https://pkg.go.dev/path/filepath#Match  

#### Parameters
`string` **`pattern`**  
Pattern to compare files with.

#### Example
```lua
--[[
	Within a folder that contains the following files:
	a.txt
	init.lua
	code.lua
	doc.pdf
]]--
local matches = fs.glob './*.lua'
print(matches)
-- -> {'init.lua', 'code.lua'}
```
</div>

<hr>
<div id='join'>
<h4 class='heading'>
fs.join(...path) -> string
<a href="#join" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Takes any list of paths and joins them based on the operating system's path separator.  

#### Parameters
`string` **`path`** (This type is variadic. You can pass an infinite amount of parameters with this type.)  
Paths to join together

#### Example
```lua
-- This prints the directory for Hilbish's config!
print(fs.join(hilbish.userDir.config, 'hilbish'))
-- -> '/home/user/.config/hilbish' on Linux
```
</div>

<hr>
<div id='mkdir'>
<h4 class='heading'>
fs.mkdir(name, recursive)
<a href="#mkdir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new directory with the provided `name`.  
With `recursive`, mkdir will create parent directories.  
-- This will create the directory foo, then create the directory bar in the  
-- foo directory. If recursive is false in this case, it will fail.  
fs.mkdir('./foo/bar', true)  

#### Parameters
`string` **`name`**  
Name of the directory

`boolean` **`recursive`**  
Whether to create parent directories for the provided name

#### Example
```lua

```
</div>

<hr>
<div id='readdir'>
<h4 class='heading'>
fs.readdir(path) -> table[string]
<a href="#readdir" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns a list of all files and directories in the provided path.  

#### Parameters
`string` **`dir`**  


</div>

<hr>
<div id='stat'>
<h4 class='heading'>
fs.stat(path) -> {}
<a href="#stat" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the information about a given `path`.  
The returned table contains the following values:  
name (string) - Name of the path  
size (number) - Size of the path in bytes  
mode (string) - Unix permission mode in an octal format string (with leading 0)  
isDir (boolean) - If the path is a directory  

#### Parameters
`string` **`path`**  


#### Example
```lua
local inspect = require 'inspect'

local stat = fs.stat '~'
print(inspect(stat))
--[[
Would print the following:
{
  isDir = true,
  mode = "0755",
  name = "username",
  size = 12288
}
]]--
```
</div>


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
|||
|----|----|
|<a href="#aliases.add">add(alias, cmd)</a>|This is an alias (ha) for the [hilbish.alias](../#alias) function.|
|<a href="#aliases.delete">delete(name)</a>|Removes an alias.|
|<a href="#aliases.list">list() -> table[string, string]</a>|Get a table of all aliases, with string keys as the alias and the value as the command.|
|<a href="#aliases.resolve">resolve(alias) -> string?</a>|Resolves an alias to its original command. Will thrown an error if the alias doesn't exist.|

<hr>
<div id='aliases.add'>
<h4 class='heading'>
hilbish.aliases.add(alias, cmd)
<a href="#aliases.add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

This is an alias (ha) for the [hilbish.alias](../#alias) function.  

#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='aliases.delete'>
<h4 class='heading'>
hilbish.aliases.delete(name)
<a href="#aliases.delete" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Removes an alias.  

#### Parameters
`string` **`name`**  


</div>

<hr>
<div id='aliases.list'>
<h4 class='heading'>
hilbish.aliases.list() -> table[string, string]
<a href="#aliases.list" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Get a table of all aliases, with string keys as the alias and the value as the command.  

#### Parameters
This function has no parameters.  
#### Example
```lua
hilbish.aliases.add('hi', 'echo hi')

local aliases = hilbish.aliases.list()
-- -> {hi = 'echo hi'}
```
</div>

<hr>
<div id='aliases.resolve'>
<h4 class='heading'>
hilbish.aliases.resolve(alias) -> string?
<a href="#aliases.resolve" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Resolves an alias to its original command. Will thrown an error if the alias doesn't exist.  

#### Parameters
`string` **`alias`**  


</div>


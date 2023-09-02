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
|<a href="#aliases.add">add(alias, cmd)</a>|This is an alias (ha) for the `hilbish.alias` function.|
|<a href="#aliases.delete">delete(name)</a>|Removes an alias.|
|<a href="#aliases.list">list() -> table<string, string></a>|Get a table of all aliases, with string keys as the alias and the value as the command.|
|<a href="#aliases.resolve">resolve(alias) -> command (string)</a>|Tries to resolve an alias to its command.|

<hr><div id='aliases.add'>
<h4 class='heading'>
hilbish.aliases.add(alias, cmd)
<a href="#aliases.add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

This is an alias (ha) for the `hilbish.alias` function.
#### Parameters
This function has no parameters.  
</div><hr><div id='aliases.delete'>
<h4 class='heading'>
hilbish.aliases.delete(name)
<a href="#aliases.delete" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Removes an alias.
#### Parameters
This function has no parameters.  
</div><hr><div id='aliases.list'>
<h4 class='heading'>
hilbish.aliases.list() -> table\<string, string>
<a href="#aliases.list" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Get a table of all aliases, with string keys as the alias and the value as the command.
#### Parameters
This function has no parameters.  
</div><hr><div id='aliases.resolve'>
<h4 class='heading'>
hilbish.aliases.resolve(alias) -> command (string)
<a href="#aliases.resolve" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Tries to resolve an alias to its command.
#### Parameters
This function has no parameters.  
</div>
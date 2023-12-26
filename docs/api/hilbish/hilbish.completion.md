---
title: Module hilbish.completion
description: tab completions
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The completions interface deals with tab completions.

## Functions
|||
|----|----|
|<a href="#completion.bins">bins(query, ctx, fields) -> entries (table), prefix (string)</a>|Return binaries/executables based on the provided parameters.|
|<a href="#completion.call">call(name, query, ctx, fields) -> completionGroups (table), prefix (string)</a>|Calls a completer function. This is mainly used to call a command completer, which will have a `name`|
|<a href="#completion.files">files(query, ctx, fields) -> entries (table), prefix (string)</a>|Returns file matches based on the provided parameters.|
|<a href="#completion.handler">handler(line, pos)</a>|This function contains the general completion handler for Hilbish. This function handles|

<hr>
<div id='completion.bins'>
<h4 class='heading'>
hilbish.completion.bins(query, ctx, fields) -> entries (table), prefix (string)
<a href="#completion.bins" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Return binaries/executables based on the provided parameters.  
This function is meant to be used as a helper in a command completion handler.  

#### Parameters
`string` **`query`**  


`string` **`ctx`**  


`table` **`fields`**  


#### Example
```lua
-- an extremely simple completer for sudo.
hilbish.complete('command.sudo', function(query, ctx, fields)
	table.remove(fields, 1)
	if #fields[1] then
		-- return commands because sudo runs a command as root..!

		local entries, pfx = hilbish.completion.bins(query, ctx, fields)
		return {
			type = 'grid',
			items = entries
		}, pfx
	end

	-- ... else suggest files or anything else ..
end)
```
</div>

<hr>
<div id='completion.call'>
<h4 class='heading'>
hilbish.completion.call(name, query, ctx, fields) -> completionGroups (table), prefix (string)
<a href="#completion.call" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Calls a completer function. This is mainly used to call a command completer, which will have a `name`  
in the form of `command.name`, example: `command.git`.  
You can check the Completions doc or `doc completions` for info on the `completionGroups` return value.  

#### Parameters
`string` **`name`**  


`string` **`query`**  


`string` **`ctx`**  


`table` **`fields`**  


</div>

<hr>
<div id='completion.files'>
<h4 class='heading'>
hilbish.completion.files(query, ctx, fields) -> entries (table), prefix (string)
<a href="#completion.files" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns file matches based on the provided parameters.  
This function is meant to be used as a helper in a command completion handler.  

#### Parameters
`string` **`query`**  


`string` **`ctx`**  


`table` **`fields`**  


</div>

<hr>
<div id='completion.handler'>
<h4 class='heading'>
hilbish.completion.handler(line, pos)
<a href="#completion.handler" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

This function contains the general completion handler for Hilbish. This function handles  
completion of everything, which includes calling other command handlers, binaries, and files.  
This function can be overriden to supply a custom handler. Note that alias resolution is required to be done in this function.  

#### Parameters
`string` **`line`**  
The current Hilbish command line

`number` **`pos`**  
Numerical position of the cursor

#### Example
```lua
-- stripped down version of the default implementation
function hilbish.completion.handler(line, pos)
	local query = fields[#fields]

	if #fields == 1 then
		-- call bins handler here
	else
		-- call command completer or files completer here
	end
end
```
</div>


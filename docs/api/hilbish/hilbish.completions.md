---
title: Module hilbish.completions
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
|<a href="#completions.call">call(name, query, ctx, fields) -> completionGroups (table), prefix (string)</a>|Calls a completer function. This is mainly used to call|
|<a href="#completions.handler">handler(line, pos)</a>|The handler function is the callback for tab completion in Hilbish.|
|<a href="#completions.bins">bins(query, ctx, fields) -> entries (table), prefix (string)</a>|Returns binary/executale completion candidates based on the provided query.|
|<a href="#completions.files">files(query, ctx, fields) -> entries (table), prefix (string)</a>|Returns file completion candidates based on the provided query.|

<hr><div id='completions.call'>
<h4 class='heading'>
hilbish.completions.call(name, query, ctx, fields) -> completionGroups (table), prefix (string)
<a href="#completions.call" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Calls a completer function. This is mainly used to call
a command completer, which will have a `name` in the form
of `command.name`, example: `command.git`.
You can check `doc completions` for info on the `completionGroups` return value.
#### Parameters
This function has no parameters.  
</div><hr><div id='completions.handler'>
<h4 class='heading'>
hilbish.completions.handler(line, pos)
<a href="#completions.handler" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

The handler function is the callback for tab completion in Hilbish.
You can check the completions doc for more info.
#### Parameters
This function has no parameters.  
</div><hr><div id='completions.bins'>
<h4 class='heading'>
hilbish.completions.bins(query, ctx, fields) -> entries (table), prefix (string)
<a href="#completions.bins" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns binary/executale completion candidates based on the provided query.
#### Parameters
This function has no parameters.  
</div><hr><div id='completions.files'>
<h4 class='heading'>
hilbish.completions.files(query, ctx, fields) -> entries (table), prefix (string)
<a href="#completions.files" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns file completion candidates based on the provided query.
#### Parameters
This function has no parameters.  
</div>
---
title: Module hilbish.abbr
description: command line abbreviations
layout: doc
menu:
  docs:
    parent: "API"
---


## Introduction
The abbr module manages Hilbish abbreviations. These are words that can be replaced
with longer command line strings when entered.
As an example, `git push` can be abbreviated to `gp`. When the user types
`gp` into the command line, after hitting space or enter, it will expand to `git push`.
Abbreviations can be used as an alternative to aliases. They are saved entirely in the history
Instead of the aliased form of the same command.

## Functions
|||
|----|----|
|<a href="#add">add(abbr, expanded|function, opts)</a>|Adds an abbreviation. The `abbr` is the abbreviation itself,|
|<a href="#remove">remove(abbr)</a>|Removes the named `abbr`.|
<hr>
<div id='remove'>
<h4 class='heading'>
hilbish.abbr.remove(abbr)
<a href="#remove" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Removes the named `abbr`.
#### Parameters
`abbr` **`string`**  


</div>

<hr>
<div id='add'>
<h4 class='heading'>
hilbish.abbr.add(abbr, expanded|function, opts)
<a href="#add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Adds an abbreviation. The `abbr` is the abbreviation itself,
while `expanded` is what the abbreviation should expand to.
It can be either a function or a string. If it is a function, it will expand to what
the function returns.
`opts` is a table that accepts 1 key: `anywhere`.
`opts.anywhere` defines whether the abbr expands anywhere in the command line or not,
whereas the default behavior is only at the beginning of the line
#### Parameters
`abbr` **`string`**  


`expanded|function` **`string`**  


`opts` **`table`**  


</div>


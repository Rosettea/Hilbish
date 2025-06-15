---
title: Module yarn
description: multi threading library
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
Yarn is a simple multithreading library. Threads are individual Lua states,
so they do NOT share the same environment as the code that runs the thread.

Example:

```lua
local yarn = require 'yarn'

-- calling t will run the yarn thread.
local t = yarn.thread(print)
t 'printing from another lua state!'
```

## Functions
|||
|----|----|
|<a href="#yarnthread">thread(fun) -> @Thread</a>|Creates a new, fresh Yarn thread.|

<hr>
<div id='yarnthread'>
<h4 class='heading'>
yarn.thread(fun) -> <a href="/Hilbish/docs/api/yarn/#thread" style="text-decoration: none;" id="lol">Thread</a>
<a href="#yarnthread" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new, fresh Yarn thread.  
`fun` is the function that will run in the thread.  

#### Parameters
This function has no parameters.  
</div>

## Types
<hr>

## Thread

### Methods

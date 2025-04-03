---
title: Module hilbish
description: No description.
layout: doc
menu:
  docs:
    parent: "Nature"
---

<hr>
<div id='runner'>
<h4 class='heading'>
hilbish.runner(mode)
<a href="#runner" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

**NOTE: This function is deprecated and will be removed in 3.0**
Use `hilbish.runner.setCurrent` instead.
This is the same as the `hilbish.runnerMode` function.
It takes a callback, which will be used to execute all interactive input.
Or a string which names the runner mode to use.
#### Parameters
`mode` **`string|function`**  


</div>

<hr>
<div id='runnerMode'>
<h4 class='heading'>
hilbish.runnerMode(mode)
<a href="#runnerMode" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

**NOTE: This function is deprecated and will be removed in 3.0**
Use `hilbish.runner.setCurrent` instead.
Sets the execution/runner mode for interactive Hilbish.
This determines whether Hilbish wll try to run input as Lua
and/or sh or only do one of either.
Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
sh, and lua. It also accepts a function, to which if it is passed one
will call it to execute user input instead.
Read [about runner mode](../features/runner-mode) for more information.
#### Parameters
`mode` **`string|function`**  


</div>


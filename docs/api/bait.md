---
title: Module bait
description: the event emitter
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

Bait is the event emitter for Hilbish. Much like Node.js and
its `events` system, many actions in Hilbish emit events.
Unlike Node.js, Hilbish events are global. So make sure to
pick a unique name!

Usage of the Bait module consists of userstanding
event-driven architecture, but it's pretty simple:
If you want to act on a certain event, you can `catch` it.
You can act on events via callback functions.

Examples of this are in the Hilbish default config!
Consider this part of it:
```lua
bait.catch('command.exit', function(code)
	running = false
	doPrompt(code ~= 0)
	doNotifyPrompt()
end)
```

What this does is, whenever the `command.exit` event is thrown,
this function will set the user prompt.

## Functions
|||
|----|----|
|<a href="#catch">catch(name, cb)</a>|Catches a hook with `name`. Runs the `cb` when it is thrown|
|<a href="#catchOnce">catchOnce(name, cb)</a>|Same as catch, but only runs the `cb` once and then removes the hook|
|<a href="#hooks">hooks(name) -> table</a>|Returns a table with hooks (callback functions) on the event with `name`.|
|<a href="#release">release(name, catcher)</a>|Removes the `catcher` for the event with `name`.|
|<a href="#throw">throw(name, ...args)</a>|Throws a hook with `name` with the provided `args`|

<hr><div id='catch'>
<h4 class='heading'>
bait.catch(name, cb)
<a href="#catch" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Catches a hook with `name`. Runs the `cb` when it is thrown
#### Parameters
`string` **`name`**  
ummm

`function` **`cb`**  
?

</div>

<hr><div id='catchOnce'>
<h4 class='heading'>
bait.catchOnce(name, cb)
<a href="#catchOnce" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Same as catch, but only runs the `cb` once and then removes the hook
#### Parameters
This function has no parameters.  
</div>

<hr><div id='hooks'>
<h4 class='heading'>
bait.hooks(name) -> table
<a href="#hooks" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns a table with hooks (callback functions) on the event with `name`.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='release'>
<h4 class='heading'>
bait.release(name, catcher)
<a href="#release" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Removes the `catcher` for the event with `name`.
For this to work, `catcher` has to be the same function used to catch
an event, like one saved to a variable.
#### Parameters
This function has no parameters.  
</div>

<hr><div id='throw'>
<h4 class='heading'>
bait.throw(name, ...args)
<a href="#throw" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Throws a hook with `name` with the provided `args`
#### Parameters
`string` **`name`**  
The name of the hook.

`any` **`args`** (This type is variadic. You can pass an infinite amount of parameters with this type.)  
The arguments to pass to the hook.

</div>


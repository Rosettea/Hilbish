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
```
bait.catch('command.exit', function(code)
	running = false
	doPrompt(code ~= 0)
	doNotifyPrompt()
end)
```

What this does is, whenever the `command.exit` event is thrown,
this function will set the user prompt.

## Functions
### bait.catch(name, cb)
Catches a hook with `name`. Runs the `cb` when it is thrown
#### Parameters
`string` **`name`**  
ummm


### bait.catchOnce(name, cb)
Same as catch, but only runs the `cb` once and then removes the hook
#### Parameters
This function has no parameters.  

### bait.hooks(name) -> table
Returns a table with hooks (callback functions) on the event with `name`.
#### Parameters
This function has no parameters.  

### bait.release(name, catcher)
Removes the `catcher` for the event with `name`.
For this to work, `catcher` has to be the same function used to catch
an event, like one saved to a variable.
#### Parameters
This function has no parameters.  

### bait.throw(name, ...args)
Throws a hook with `name` with the provided `args`
#### Parameters
`string` **`name`**  
The name of the hook.

`any` **`args`** (This type is variadic. You can pass an infinite amount of parameters with this type.)  
The arguments to pass to the hook.



---
title: Module bait
description: the event emitter
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
Bait is the event emitter for Hilbish. Why name it bait? Why not.
It throws hooks that you can catch. This is what you will use if
you want to listen in on hooks to know when certain things have
happened, like when you've changed directory, a command has failed,
etc. To find all available hooks thrown by Hilbish, see doc hooks.

## Functions
### bait.catch(name, cb)
Catches a hook with `name`. Runs the `cb` when it is thrown
#### Parameters
This function has no parameters.  

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



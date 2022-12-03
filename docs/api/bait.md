---
name: Module bait
description: the event emitter
layout: apidoc
---

## Introduction
Bait is the event emitter for Hilbish. Why name it bait? Why not.
It throws hooks that you can catch. This is what you will use if
you want to listen in on hooks to know when certain things have
happened, like when you've changed directory, a command has failed,
etc. To find all available hooks thrown by Hilbish, see doc hooks.

## Functions
### catchOnce(name, cb)
Same as catch, but only runs the `cb` once and then removes the hook

### hooks(name) -> {cb, cb...}
Returns a table with hooks on the event with `name`.

### release(name, catcher)
Removes the `catcher` for the event with `name`
For this to work, `catcher` has to be the same function used to catch
an event, like one saved to a variable.

### throw(name, ...args)
Throws a hook with `name` with the provided `args`

### catch(name, cb)
Catches a hook with `name`. Runs the `cb` when it is thrown


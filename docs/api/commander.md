---
title: Module commander
description: library for custom commands
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

Commander is a library for writing custom commands in Lua.
In order to make it easier to write commands for Hilbish,
not require separate scripts and to be able to use in a config,
the Commander library exists. This is like a very simple wrapper
that works with Hilbish for writing commands. Example:

```lua
local commander = require 'commander'

commander.register('hello', function(args, sinks)
	sinks.out:writeln 'Hello world!'
end)
```

In this example, a command with the name of `hello` is created
that will print `Hello world!` to output. One question you may
have is: What is the `sinks` parameter?

The `sinks` parameter is a table with 3 keys: `in`, `out`,
and `err`. The values of these is a <a href="/Hilbish/docs/api/hilbish/#sink" style="text-decoration: none;">Sink</a>.

- `in` is the standard input. You can read from this sink
to get user input. (**This is currently unimplemented.**)
- `out` is standard output. This is usually where text meant for
output should go.
- `err` is standard error. This sink is for writing errors, as the
name would suggest.

## Functions
### commander.deregister(name)
Deregisters any command registered with `name`
#### Parameters
This function has no parameters.  

### commander.register(name, cb)
Register a command with `name` that runs `cb` when ran
#### Parameters
This function has no parameters.  


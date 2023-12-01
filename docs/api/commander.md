---
title: Module commander
description: library for custom commands
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

Commander is the library which handles Hilbish commands. This makes
the user able to add Lua-written commands to their shell without making
a separate script in a bin folder. Instead, you may simply use the Commander
library in your Hilbish config.

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
and `err`. All of them are a <a href="/Hilbish/docs/api/hilbish/#sink" style="text-decoration: none;">Sink</a>.

- `in` is the standard input.
You may use the read functions on this sink to get input from the user.
- `out` is standard output.
This is usually where command output should go.
- `err` is standard error.
This sink is for writing errors, as the name would suggest.

## Functions
|||
|----|----|
|<a href="#deregister">deregister(name)</a>|Removes the named command. Note that this will only remove Commander-registered commands.|
|<a href="#register">register(name, cb)</a>|Adds a new command with the given `name`. When Hilbish has to run a command with a name,|

<hr><div id='deregister'>
<h4 class='heading'>
commander.deregister(name)
<a href="#deregister" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Removes the named command. Note that this will only remove Commander-registered commands.
#### Parameters
`string` **`name`**  
Name of the command to remove.

</div>

<hr><div id='register'>
<h4 class='heading'>
commander.register(name, cb)
<a href="#register" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Adds a new command with the given `name`. When Hilbish has to run a command with a name,
it will run the function providing the arguments and sinks.


#### Parameters
`string` **`name`**  
Name of the command

`function` **`cb`**  
Callback to handle command invocation

#### Example
```lua
-- When you run the command `hello` in the shell, it will print `Hello world`.
-- If you run it with, for example, `hello Hilbish`, it will print 'Hello Hilbish'
commander.register('hello', function(args, sinks)
	local name = 'world'
	if #args > 0 then name = args[1] end

	sinks.out:writeln('Hello ' .. name)
end)
````
</div>


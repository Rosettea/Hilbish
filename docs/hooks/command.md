---
title: Command
description:
layout: doc
menu:
  docs:
    parent: "Signals"
---

## command.preexec
Thrown right before a command is executed.

#### Variables
`string` **`input`**  
The raw string that the user typed. This will include the text
without changes applied to it (argument substitution, alias expansion,
etc.)

`string` **`cmdStr`**  
The command that will be directly executed by the current runner.

<hr>

## command.exit
Thrown after the user's ran command is finished.

#### Variables
`number` **`code`**  
The exit code of what was executed.

`string` **`cmdStr`**  
The command or code that was executed

<hr>
	
## command.not-found
Thrown if the command attempted to execute was not found.
This can be used to customize the text printed when a command is not found.
Example:
```lua
local bait = require 'bait'
-- Remove any present handlers on `command.not-found`

local notFoundHooks = bait.hooks 'command.not-found'
for _, hook in ipairs(notFoundHooks) do
	bait.release('command.not-found', hook)
end

-- then assign custom
bait.catch('command.not-found', function(cmd)
	print(string.format('The command "%s" was not found.', cmd))
end)
```

#### Variables
`string` **`cmdStr`**  
The name of the command.

<hr>
	
## command.not-executable
Thrown when the user attempts to run a file that is not executable
(like a text file, or Unix binary without +x permission).

#### Variables
`string` **`cmdStr`**  
The name of the command.

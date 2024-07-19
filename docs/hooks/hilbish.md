---
title: Hilbish
description:
layout: doc
menu:
  docs:
    parent: "Signals"
---

## hilbish.exit
Sent when Hilbish is going to exit.

#### Variables
This signal returns no variables.

<hr>

## hilbish.vimMode
Sent when the Vim mode of Hilbish is changed (like from insert to normal mode).
This can be used to change the prompt and notify based on Vim mode.

#### Variables
`string` **`modeName`**  
The mode that has been set.
Can be these values: `insert`, `normal`, `delete` or `replace`

<hr>

## hilbish.cancel
Sent when the user cancels their command input with Ctrl-C

#### Variables
This signal returns no variables.

<hr>

## hilbish.notification
Thrown when a [notification](../../features/notifications) is sent.

#### Variables
`table` **`notification`**  
The notification. The properties are defined in the link above.

<hr>

+ `hilbish.vimAction` -> actionName, args > Sent when the user does a "vim action," being something
like yanking or pasting text. See `doc vim-mode actions` for more info.

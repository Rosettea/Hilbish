---
title: Hilbish
description:
layout: doc
menu:
  docs:
    parent: "Signals"
---

+ `hilbish.exit` > Sent when Hilbish is about to exit.

+ `hilbish.vimMode` -> modeName > Sent when Hilbish's Vim mode is changed (example insert to normal mode),
`modeName` is the name of the mode changed to (can be `insert`, `normal`, `delete` or `replace`).

+ `hilbish.vimAction` -> actionName, args > Sent when the user does a "vim action," being something
like yanking or pasting text. See `doc vim-mode actions` for more info.

+ `hilbish.cancel` > Sent when the user cancels their input with Ctrl-C.

+ `hilbish.notification` -> message > Sent when a message is
sent.

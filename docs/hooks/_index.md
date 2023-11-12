---
title: Hooks
description:
layout: doc
weight: -50
menu:
  docs
---

Here is a list of bait hooks that are thrown by Hilbish. If a hook is related
to a command, it will have the `command` scope, as example.

Here is the format for a doc for a hook:
+ {hook name} -> args > description

`{args}` just means the arguments of the hook. If a hook doc has the format
of `arg...`, it means the hook can take/recieve any number of `arg`.

+ error -> eventName, handler, err > Emitted when there is an error in
an event handler. The `eventName` is the name of the event the handler
is for, the `handler` is the callback function, and `err` is the error
message.

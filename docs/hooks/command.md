---
title: Command
description:
layout: doc
menu:
  docs:
    parent: "Hooks"
---

- `command.preexec` -> input, cmdStr > Thrown before a command
is executed. The `input` is the user written command, while `cmdStr`
is what will be executed (`input` will have aliases while `cmdStr`
will have alias resolved input).

- `command.exit` -> code, cmdStr > Thrown when a command exits.
`code` is the exit code of the command, and `cmdStr` is the command that was run.

- `command.not-found` -> cmdStr > Thrown when a command is not found.

- `command.not-executable` -> cmdStr > Thrown when Hilbish attempts to run a file
that is not executable.

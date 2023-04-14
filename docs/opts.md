---
title: "Opts"
date: 2023-04-14T01:01:10-04:00
draft: false
---

Opts are simple settings and switches to control certain Hilbish behavior.
This ranges from things like the greeting messages on bootup to "autocd"
functionality.

They can be changed via a simple assign (like `hilbish.opts.namehere = true`).

# Available Opts
## `autocd`
The `autocd` opt (default false) makes it so that if a lone path is
ran, Hilbish will change directory to that path.

Example:  
```
sammy ~ ∆ hilbish.opts.autocd = true
sammy ~ ∆ ~/Downloads
sammy ~/Downloads ∆ 
```

## `history`
This opt controls if Hilbish will store commands in history. It
is default true.

## `greeting`
The `greeting` is exactly as the name says: a greeting message on Hilbish
startup. It can be set to any value, and when Hilbish finishes initalizing
will print that text. The `greeting` will be passed to Lunacolors for
formatting.

## `motd`
The `motd` is a message to summarize the current running version. It can be
set as a boolean (true or false).

---
title: Runner Mode
description: Customize the interactive script/command runner.
layout: doc
menu: 
  docs:
    parent: "Features"
---

Hilbish allows you to change how interactive text can be interpreted.
This is mainly due to the fact that the default method Hilbish uses
is that it runs Lua first and then falls back to shell script.

In some cases, someone might want to switch to just shell script to avoid
it while interactive but still have a Lua config, or go full Lua to use
Hilbish as a REPL. This also allows users to add alternative languages,
instead of either like Fennel.

---
title: Options
description: Simple customizable options.
layout: doc
menu: 
  docs:
    parent: "Features"
---

Opts are simple toggle or value options a user can set in Hilbish.
As toggles, there are things like `autocd` or history saving. As values,
there is the `motd` which the user can either change to a custom string or disable.

Opts are accessed from the `hilbish.opts` table. Here they can either
be read or modified

### `autocd`
#### Value: `boolean`
#### Default: `false`

The autocd opt makes it so that lone directories attempted to be executed are
instead set as the shell's directory.

Example:
```
~/Directory                                     
∆ ~
~                                                                             
∆ Downloads
~/Downloads                                                                   
∆ ../Documents
~/Documents                                                                   
∆ 
```

<hr>

### `history`
#### Value: `boolean`
#### Default: `true`
Sets whether command history will be saved or not.

<hr>
	
### `greeting`
#### Value: `boolean` or `string`
The greeting is the message that Hilbish shows on startup
(the one which says Welcome to Hilbish).

This can be set to either true/false to enable/disable or a custom greeting string.

<hr>

### `motd`
#### Value: `boolean`
#### Default: `true`
The message of the day shows the current major.minor version and
includes a small range of things added in the current release.

This can be set to `false` to disable the message.

<hr>

### `fuzzy`
#### Value: `boolean`
#### Default: `false`
Toggles the functionality of fuzzy history searching, usable
via the menu in Ctrl-R. Fuzzy searching is an approximate searching
method, which means results that match *closest* will be shown instead
of an exact match.

<hr>

### `notifyJobFinish`
#### Value: `boolean`
#### Default: `true`
If this is enabled, when a background job is finished,
a [notification](../notifications) will be sent.

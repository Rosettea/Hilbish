---
title: Getting Started
layout: doc
weight: -10
---

To start Hilbish, open a terminal. If Hilbish has been installed and is not the
default shell, you can simply run `hilbish` to start it. This will launch
a normal interactive session.
To exit, you can either run the `exit` command or hit Ctrl+D.

# Setting as Default
## Login shell
There are a few ways to make Hilbish your default shell. A simple way is 
to make it your user/login shell.

{{< warning `It is not recommended to set Hilbish as your login shell. That is expected to be a 
POSIX compliant shell, which Hilbish is not. At most, there will just be a 
few variables missing in your environment` >}}

To do that, simply run `chsh -s /usr/bin/hilbish`.
Some distros (namely Fedora) might have `lchsh` instead, which is used like `lchsh <user>`.
When prompted, you can put the path for Hilbish.

## Default with terminal
The simpler way is to set the default shell for your terminal. The way of 
doing this depends on how your terminal settings are configured.

## Run after login shell
Some shells (like zsh) have an rc file, like `.zlogin`, which is ran when the shell session
is a login shell. In that file, you can run Hilbish. Example:

```
exec hilbish -S -l
```

This will replace the shell with Hilbish, set $SHELL to Hilbish and launch it as a login shell.

# Configuration
Once installation and setup has been done, you can then configure Hilbish.
It is configured and scripted via Lua, so the config file is a Lua file.
You can use any pure Lua library to do whatever you want.

Hilbish's sample configuration is usually located in `hilbish.dataDir .. '/.hilbishrc.lua'`.
You can print that path via Lua to see what it is: `print(hilbish.dataDir .. '/.hilbishrc.lua')`.
As an example, it will usually will result in `/usr/share/hilbish/.hilbishrc.lua` on Linux.

To edit your user configuration, you can copy that file to `hilbish.userDir.config .. '/hilbish/init.lua'`,
which follows XDG on Linux and MacOS, and is located in %APPDATA% on Windows.

As the directory is usually `~/.config` on Linux, you can run this command to copy it:  
`cp /usr/share/hilbish/.hilbishrc.lua ~/.config/hilbish/init.lua`

Now you can get to editing it. Since it's just a Lua file, having basic
knowledge of Lua would help. All of Lua's standard libraries and functions
from Lua 5.4 are available. Hilbish has some custom and modules that are
available. To see them, you can run the `doc` command. This also works as
general documentation for other things.

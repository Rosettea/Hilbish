# 🎀 Changelog

This is the changelog for the Hilbish shell made in Go and Lua.

## [0.5.1] - 2021-06-16

## Added

- Add `~/.config/hilbish` as a require path

## Changed

- `cd` hook is only thrown after directory has actually changed

## Fixed

- Handle error in commander properly, preventing a panic from Lua

## [0.5.0] - 2021-06-12

An absolutely massive release. Probably the biggest yet, includes a bunch of fixes and new features and convenient additions to the Lua API.

## Added

- `-n` flag, which checks Lua for syntax errors without running it
- `exec(command)` function, acts like the `exec` builtin in sh 
    - Example: `exec 'awesome'` in an .xinitrc file with Hilbish as shebang
- Commands from commander can now `return` an exit code 
```lua
commander.register('false', function()
return 1
end)
```
    - When `false` is run, it will have the exit code of `1`, this is shorter/easier than throwing the command.exit hook and can work if the functionality of that changes
- Added `-c` description
- `args` variable, set when Hilbish runs a Lua script. It is an array that includes the execute path as the first argument
- Lua code can be aliased
- Recursive aliases
    - At the moment this only works for the first argument
- Hilbish can now be used with Hilbiline if compiled to do so (currently only for testing purposes)
- `goro(func)` runs a `func`tion in a goroutine. With channels that gopher-lua also provides, one can do parallelism and concurrency in Lua (but go style). 
    - `coroutine` no those dont exist they dont matter `goro` is easier
- `cd -` will change to the previous directory
- `hilbish.cwd()` gets the current working directory
- `timeout(func, time)` works exactly like the `setTimeout` function in JavaScript. It will run `func` after a period of `time` in milliseconds.
- `interval(func, time)` works exactly like the `setInterval` function in JavaScripit. It will run `func` every `time` milliseconds
- `hilbish.home` is a crossplatform Lua alternative to get the home directory easily.
- `commander.deregister(cmdName)` de-registers any command defined with commander.

## Changed

- **Breaking Change**: Move `_user` and `_ver` to a global `hilbish` table
    - Accessing username and Hilbish version is now done with `hilbish.user` and `hilbish.ver`
- `hilbish.run(cmd)` runs cmd with Hilbish's sh interpreter. Using this function instead of `os.execute` ensures that sh syntax works everywhere Hilbish does.
- `hilbish.flag(flag)` checks if flag has been passed to Hilbish.
- Aliases now work with every command and not only the first one
    - Therefore `alias1; alias2` works now
- `command.not-found` hook
- `$SHLVL` is now incremented in Hilbish. If not a valid number, it will be changed to 1
- `fs.mkdir` can now make directories recursively if the 2nd argument is set to `true`
    - `fs.mkdir('path/to/dir', true)`
- Hilbish runs a `preload.lua` file in the current directory first, then falls back to the global preload. Before the order was reversed.
- Check if aliased command is defined in Lua, so registered `commander`s can be aliased
- Add input to history before alias expansion. Basically, this adds the actual alias to history instead of the aliased command.
- Global preload path, require paths, default config directory and sample config directory can now be changed at compile time to help support other systems.

## Fixed

- `cd` now exits with code `1` instead of the error code if it occurs
- Don't append directory to $PATH with `appendPath` if its already there
- Continued input is no longer joined with a space unless explicitly wanted
- Hilbish won't try to go interactive if it isn't launched in a TTY (terminal)
- Ctrl+D on a continue prompt with no input no longer causes a panic
- Actually handle the `-h`/`--help` option

## [0.4.0] - 2021-05-01

## Added
- Ctrl C in the prompt now cancels/clear input (I've needed this for so long also)
- Made Hilbish act like a login shell on login 
    - If Hilbish is the login shell, or the `-l`/`--login` flags are used, Hilbish will use an additional `~/.hprofile.lua` file, you can use this to set environment variables once on login
- `-c` has been added to run a single command (this works exactly like being in the prompt would, so Lua works as well)
- `-i` (also `--interactive`) has been added to force Hilbish to be an interactive shell in cases where it usually wont be (like with `-c`)
- Use readline in continue prompt
- Added a `mulitline` hook that's thrown when in the continue/multiline prompt
- Added `appendPath` function to append a directory to `$PATH`
    - `~` will be expanded to `$HOME` as well
- A utility `string.split` function is now added 
    - `string.split(str, delimiter)`
- Added a `_user` variable to easily get current user's name

## Changed

- **BREAKING Change**: [Lunacolors](https://github.com/Hilbis/Lunacolors) has replaced ansikit for formatting colors, which means the format function has been removed from ansikit and moved to Lunacolors. 
    - Users must replace ansikit with `lunacolors` in their config files
- A getopt-like library is now used for command line flag parsing
- `cd` builtin now supports using environment variables
    - This means you can now `cd $NVM_DIR` as an example
- Function arguments are now more strictly typed (`prompt(nil)` wouldnt work now)
- Other general code/style changes

## Fixed

- Fix makefile adding Hilbish to `/etc/shells` on every `make install`

Since Lunacolors is a submodule, you may just want to completely reclone Hilbish recursively and then update (rerun `make install`)
Or instead of recloning, run `git submodule update --init --recursive` in Hilbish's git directory

## [0.3.2] - 2021-04-10

### Added

- Add more functions to `ansikit` module
- Add functions `stat` and `mkdir` to `fs` module
- `-C` flag to define path to config
- Add require path `~/.local/share/hilbish/libs`

### Changed

- Continue to new line if output doesnt end with it

Observed:
![Observed](https://camo.githubusercontent.com/5be15fed950a2926e6f14dfe4427b84b7c0c448d5d937f9df15959ca934a50ce/68747470733a2f2f6d6f646575732e69732d696e736964652e6d652f70633335416133492e706e67)

## [0.3.1] - 2021-04-06

### Fixed

- Fix `%u` in prompt format being full name and instead make it actually username

## [0.3.0] - 2021-04-05

### Added

- Added a `multiprompt` function to change the prompt of the multiline/continuation/newline prompt
- `_ver` variable to get Hilbish's version from Lua

### Changed

- **BREAKING Change**: Removed Bait hooks `command.success` and `command.fail`, there is now the single hook `command.exit`, with a single argument passed which the exit code of the command. Use this to determine if a command has failed or not (failure is code != 0)
- **BREAKING Change**: The Ansikit function `text` has been renamed to `format`.
Input ending 
- `fs.cd` now throws an exception instead of silently failing, which you should handle with `pcall`
- Enhancements to the `cd` command:
    - With no arguments will move to $HOME
    - Now throws a cd hook, with a single hook arg being the arguments to the command
    - Now works for directories with spaces
- Lua input now throws a success hook if it succeeded
- Command history is now saved to `~/.hilbish-history`
- Globals defined in Lua that are strings will be defined as an env variable ([#16](https://github.com/Rosettea/Hilbish/pull/16))
- Input ending with `\` will now go to a newline
- `exit` command is now written in Lua

### Fixed

- Input is now trimmed
- Internals are slightly cleaned up and codebase split up
- Hilbish will now fall back to a builtin minimal config if the user's config has syntax errors on loading
- Commands defined in Lua now report the error to the user cleanly instead of panicking if it has occured

## [0.2.0] - 2021-03-31

### Added

- Hooks (events) are the new and main thing in v0.2, you can now listen for hooks or cast out (emit) custom ones, via the [bait](https://github.com/Hilbis/Hilbish/wiki/Bait) package
- `^^` to refer to the previous command. It's for the lazy hands like me, so I can do something like `ENV=VAR ^^`
- Way more (more like any) comments in the core code.

### Changed

- Prompt has been changed to have printf-like verbs to format. This makes it easier on the user's side to configure, and doesn't require hooks to change it for things like current directory.
- Default prompt's changed and the triangle changes color based on if command failed or not.

## [0.1.2] - 2021-03-24

### Added

- Add Lua input to history

## [0.1.1] - 2021-03-24

### Added

- Go to new line if sh input is incomplete

```bash
> for i in {1..5}
```

This input for example will prompt for more input to complete:

![input](https://camo.githubusercontent.com/b757e474da5880d57be135087f59f45ab214b8f39f182b299d861cac7b6d84ff/68747470733a2f2f6d6f646575732e69732d696e736964652e6d652f30624456547461352e706e67)

## [0.1.0] - 2021-03-24

### Added

- Tab complete files
- Makefile installation
- sh support

## [0.0.12] - 2021-03-21

First "stable" release of Hilbish.

[0.5.1]: https://github.com/Rosettea/Hilbish/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/Rosettea/Hilbish/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/Rosettea/Hilbish/compare/v0.3.2...v0.4.0
[0.3.2]: https://github.com/Rosettea/Hilbish/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/Rosettea/Hilbish/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/Rosettea/Hilbish/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/Rosettea/Hilbish/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/Rosettea/Hilbish/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/Rosettea/Hilbish/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/Rosettea/Hilbish/compare/v0.0.12...v0.1.0 
[0.0.12]: https://github.com/Rosettea/Hilbish/releases/tag/v0.0.12

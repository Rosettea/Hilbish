# 🎀 Changelog

## [1.0.4] - 2021-03-12
### Fixed
- Panic when history directory doesn't exist

## [1.0.3] - 2021-03-12
### Fixed
- Removed duplicate executable suggestions
- User input is added to history now instead of what's ran by Hilbish
- Formatting issue with prompt on no input

## [1.0.2] - 2021-03-06
### Fixed
- Cases where Hilbish's history directory doesn't exist will no longer cause a panic

## [1.0.1] - 2021-03-06
### Fixed
- Using `hilbish.appendPath` will no longer result in string spam (debugging thing left being)
- Prompt gets set properly on startup

## [1.0.0] - 2021-03-06
### Added
- MacOS is now officialy supported, default compile time vars have been added
for it
- Windows is properly supported as well
- `catchOnce()` to bait - catches a hook once
- `hilbish.aliases` interface - allows you to add, delete and list all aliases
with Lua
- `hilbish.appendPath()` can now take a table of arguments for ease of use
- `hilbish.which(binName)` acts as the which builtin for other shells,
it finds the path to `binName` in $PATH
- Signal hooks `sigusr1` and `sigusr2` (unavailable on Windows)
- Commands starting with a space won't be added to history
- Vim input mode
  - Hilbish's input mode for text can now be changed to either Emacs
  (like it always was) or Vim via `hilbish.inputMode()`
  - Changing Vim mode throws a `hilbish.vimMode` hook
  - The current Vim mode is also accessible with the `hilbish.vimMode` property
- Print errors in `hilbish.timeout()` and `hilbish.goro()` callbacks
- `hilbish.exit` hook is thrown when Hilbish is going to exit
- `hilbish.exitCode` property to get the exit code of the last executed command
- `screenMain` and `screenAlt` functions have been added to Ansikit to switch
to the terminal's main and alt buffer respectively

### Fixed
- Tab completion for executables
- Stop interval (`hilbish.interval()`) when an error occurs
- Errors in bait hooks no longer cause a panic, and remove the handler for the hook as well
- Formatting of home dir to ~
- Check if Hilbish is in interactive before trying to use its handlers for signals
- Global `args` table when running as script is no longer userdata
- Home dir is now added to recent dirs (the case of cd with no arg)
- `index` subdoc will no longer appear
- Alias expansion with quotes
- Add full command to history in the case of incomplete input
- `hilbish.exec()` now has a windows substitute
- Fixed case of successful command after prompted for more input not writing to history
- `command.exit` is thrown when sh input is incorrect and when command executed after continue
prompt exits successfully

### Changed
- The minimal config is truly minimal now
- Default config is no longer copied to user's config and is instead ran its location
#### Breaking Changes
(there were a lot...)
- Change default SHLVL to 0 instead of 1
- ~/.hilbishrc.lua will no longer be run by default, it now
only uses the paths mentioned below.
- Changed Hilbish's config path to something more suited
according to the OS (`$XDG_CONFIG_HOME/hilbish/init.lua` on Linux,
`~/Library/Application Support/hilbish/init.lua` on MacOS and
(`%APPDATA%/hilbish/init.lua` on Windows). Previously on Unix-like it was
`$XDG_CONFIG_HOME/hilbish/hilbishrc.lua`
- The history path has been changed to a better suited path.
On Linux, it is `$XDG_DATA_HOME/hilbish/.hilbish-history` and for others it is
the config path.
- `hilbish.xdg` no longer exists, use `hilbish.userDir` instead,
as it functions the same but is OS agnostic
- `hilbish.flag()` has been removed
- `~/.hprofile.lua` has been removed, instead check in your config if `hilbish.login`
is true
- `hilbish.complete()` has had a slight refactor to fit with the new readline library.
It now expects a table of "completion groups" which are just tables with the
`type` and `items` keys. Here is a (more or less) complete example of how it works now:
	```lua
	hilbish.complete('command.git', function()
		return {
			{
				items = {
					'add',
					'clone'
				},
				type = 'grid'
			},
			{
				items = {
					['--git-dir'] = {'Description of flag'},
					'-c'
				},
				type = 'list'
			}
		}
	end)
	```
	Completer functions are now also expected to handle subcommands/subcompletions

## [0.7.1] - 2021-11-22
### Fixed
- Tab complete absolute paths to binaries properly
- Allow execution of absolute paths to binaries (https://github.com/Rosettea/Hilbish/commit/06272778f85dad04e0e7abffc78a5b9b0cebd067 regression)

## [0.7.0] - 2021-11-22
### Added
- `hilbish.interactive` and `hilbish.login` properties to figure out if Hilbish is interactive or a login shell, respectively.
- `hilbish.read` function to take input more elegantly than Lua's `io.read`
- Tab Completion Enhancements
  - A new tab complete API has been added. It is the single `complete` function which takes a "scope" (example: `command.<cmdname>`) and a callback which is
  expected to return a table. Users can now add custom completions for specific commands.
  An example is:
  ```lua
  complete('command.git', function()
	return {
		'add',
		'version',
		commit = {
			'--message',
			'--verbose',
			'<file>'
		}
	}
  end)
  ```
  For `git`, Hilbish will complete commands add, version and commit. For the commit subcommand, it will complete the flags and/or files which `<file>` is used to represent.
  - Hilbish will now complete binaries in $PATH, or any executable to a path (like `./` or `../`)
  - Files with spaces will be automatically put in quotes and completions will work for them now.
- `prependPath` function (#81)
- Signal hooks (#80)
  - This allows scripts to add their own way of handling terminal resizes (if you'd need that) or Ctrl-C
- Module properties (like `hilbish.ver`) are documented with the `doc` command.
- Document bait hooks

### Fixed
- The prompt won't come up on terminal resize anymore.
- `appendPath` should work properly on Windows.
- A panic when a commander has an error has been fixed.

## [0.6.1] - 2021-10-21
### Fixed
- Require paths now use the `dataDir` variable so there is no need to change it anymore unless you want to add more paths
- Remove double slash in XDG data require paths
- Ctrl+C is handled properly when not interactive and won't result in a panic anymore
- Commanders are handled by the sh interpreter library only now, so they work with sh syntax

### Changed
- Error messages from `fs` functions now include the path provided

## [0.6.0] - 2021-10-17
### Added
- Hilbish will expand `~` in the preloadPath and samplePathConf variables. These are for compile time.
- On Windows, the hostname in `%u` has been removed.
- Made it easier to compile on Windows by adding Windows-tailored vars and paths.
- Add require paths `./libs/?/?.lua`
- Hilbish will now respect $XDG_CONFIG_HOME and will load its config and history there first and use Lua libraries in there and $XDG_DATA_HOME if they are set. (#71)
  - If not, Hilbish will still default to `~`
- Added some new hooks
  - `command.precmd` is thrown right before Hilbish prompts for input
  - `command.preexec` is thrown right before Hilbish executes a command. It passes 2 arguments: the command as the user typed, and what Hilbish will actually execute (resolved alias)
- `hilbish.dataDir` is now available to know the directory of Hilbish data files (default config, docs, preload, etc)
- A `docgen` program has been added to `cmd/docgen` in the GitHub repository, As the name suggests, it will output docs in a `docs` folder for functions implemented in Go
- All hilbish modules/libraries now have a `__doc` metatable entry which is simply a short description of the module.
- `fs.readdir(dir)` has been added. It will return a table of files in `dir`
- Errors in the `fs.mkdir` function are now handled.
- **Breaking Change:** `fs.cd` no longer returns a numeric code to indicate error. Instead, it returns an error message.
- The `doc` command has been added to document functions of Hilbish libraries. Run the command for more details.
- `link(url, text)` has been added to `ansikit`. It returns a string which can be printed to produce a hyperlink in a terminal. Note that not all terminals support this feature.
- The [Succulent](https://github.com/Rosettea/Succulent) library has been added. This includes more utility functions and expansions to the Lua standard library itself.
- The command string is now passed to the `command.exit` hook

### Changed
- Hilbish won't print an extra newline at exit with ctrl + d
- `command.exit` with 0 exit code will now be thrown if input is nothing
- **Breaking Change:** `fs.stat` has been made better. It returns a proper table instead of userdata, and has fields instead of functions
  - It includes `name`, `mode` as a octal representation in a string, `isDir`, and `size`

### Fixed
- `timeout()` is now blocking
- Directories with spaces in them can now be `cd`'d to
- An alias with the same name as the command will now not cause a freeze (#73)
- Userdata is no longer returned in the following cases:
  - Commander arguments
  - `fs` functions

## [0.5.1] - 2021-06-16

### Added

- Add `~/.config/hilbish` as a require path

### Changed

- `cd` hook is only thrown after directory has actually changed

### Fixed

- Handle error in commander properly, preventing a panic from Lua

## [0.5.0] - 2021-06-12

An absolutely massive release. Probably the biggest yet, includes a bunch of fixes and new features and convenient additions to the Lua API.

### Added

- `-n` flag, which checks Lua for syntax errors without running it
- `exec(command)` function, acts like the `exec` builtin in sh
    - Example: `exec 'awesome'` in an .xinitrc file with Hilbish as shebang
- Commands from commander can now `return` an exit code
```lua
commander.register('false', function()
return 1
end)
```
When `false` is run, it will have the exit code of `1`, this is shorter/easier than throwing the command.exit hook and can work if the functionality of that changes
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

### Changed

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

### Fixed

- `cd` now exits with code `1` instead of the error code if it occurs
- Don't append directory to $PATH with `appendPath` if its already there
- Continued input is no longer joined with a space unless explicitly wanted
- Hilbish won't try to go interactive if it isn't launched in a TTY (terminal)
- Ctrl+D on a continue prompt with no input no longer causes a panic
- Actually handle the `-h`/`--help` option

## [0.4.0] - 2021-05-01

### Added
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

### Changed

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

[1.0.4]: https://github.com/Rosettea/Hilbish/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/Rosettea/Hilbish/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/Rosettea/Hilbish/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/Rosettea/Hilbish/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/Rosettea/Hilbish/compare/v0.7.1...v1.0.0
[0.7.1]: https://github.com/Rosettea/Hilbish/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/Rosettea/Hilbish/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/Rosettea/Hilbish/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/Rosettea/Hilbish/compare/v0.5.1...v0.6.0
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

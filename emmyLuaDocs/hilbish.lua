--- @meta

local hilbish = {}

--- This is an alias (ha) for the [hilbish.alias](../#alias) function.
--- @param alias string
--- @param cmd string
function hilbish.aliases.add(alias, cmd) end

--- This is the same as the `hilbish.runnerMode` function.
--- It takes a callback, which will be used to execute all interactive input.
--- In normal cases, neither callbacks should be overrided by the user,
--- as the higher level functions listed below this will handle it.
function hilbish.runner.setMode(cb) end

--- Deletes characters in the line by the given amount.
function hilbish.editor.deleteByAmount(amount) end

--- Returns the current input line.
function hilbish.editor.getLine() end

--- Returns the text that is at the register.
function hilbish.editor.getVimRegister(register) end

--- Inserts text into the Hilbish command line.
function hilbish.editor.insert(text) end

--- Reads a keystroke from the user. This is in a format of something like Ctrl-L.
function hilbish.editor.getChar() end

--- Sets the vim register at `register` to hold the passed text.
function hilbish.editor.setVimRegister(register, text) end

--- Return binaries/executables based on the provided parameters.
--- This function is meant to be used as a helper in a command completion handler.
--- 
--- 
function hilbish.completion.bins(query, ctx, fields) end

--- Calls a completer function. This is mainly used to call a command completer, which will have a `name`
--- in the form of `command.name`, example: `command.git`.
--- You can check the Completions doc or `doc completions` for info on the `completionGroups` return value.
function hilbish.completion.call(name, query, ctx, fields) end

--- Returns file matches based on the provided parameters.
--- This function is meant to be used as a helper in a command completion handler.
function hilbish.completion.files(query, ctx, fields) end

--- This function contains the general completion handler for Hilbish. This function handles
--- completion of everything, which includes calling other command handlers, binaries, and files.
--- This function can be overriden to supply a custom handler. Note that alias resolution is required to be done in this function.
--- 
--- 
function hilbish.completion.handler(line, pos) end

--- Sets an alias, with a name of `cmd` to another command.
--- 
--- 
function hilbish.alias(cmd, orig) end

--- Appends the provided dir to the command path (`$PATH`)
--- 
--- 
function hilbish.appendPath(dir) end

--- Registers a completion handler for the specified scope.
--- A `scope` is expected to be `command.<cmd>`,
--- replacing <cmd> with the name of the command (for example `command.git`).
--- The documentation for completions, under Features/Completions or `doc completions`
--- provides more details.
--- 
--- 
function hilbish.complete(scope, cb) end

--- Returns the current directory of the shell.
function hilbish.cwd() end

--- Replaces the currently running Hilbish instance with the supplied command.
--- This can be used to do an in-place restart.
function hilbish.exec(cmd) end

--- Puts `fn` in a Goroutine.
--- This can be used to run any function in another thread at the same time as other Lua code.
--- **NOTE: THIS FUNCTION MAY CRASH HILBISH IF OUTSIDE VARIABLES ARE ACCESSED.**
--- **This is a limitation of the Lua runtime.**
function hilbish.goro(fn) end

--- Line highlighter handler.
--- This is mainly for syntax highlighting, but in reality could set the input
--- of the prompt to *display* anything. The callback is passed the current line
--- and is expected to return a line that will be used as the input display.
--- Note that to set a highlighter, one has to override this function.
--- 
function hilbish.highlighter(line) end

--- The command line hint handler. It gets called on every key insert to
--- determine what text to use as an inline hint. It is passed the current
--- line and cursor position. It is expected to return a string which is used
--- as the text for the hint. This is by default a shim. To set hints,
--- override this function with your custom handler.
--- 
--- 
function hilbish.hinter(line, pos) end

--- Sets the input mode for Hilbish's line reader.
--- `emacs` is the default. Setting it to `vim` changes behavior of input to be
--- Vim-like with modes and Vim keybinds.
function hilbish.inputMode(mode) end

--- Runs the `cb` function every specified amount of `time`.
--- This creates a timer that ticking immediately.
function hilbish.interval(cb, time) end

--- Changes the text prompt when Hilbish asks for more input.
--- This will show up when text is incomplete, like a missing quote
--- 
--- 
function hilbish.multiprompt(str) end

--- Prepends `dir` to $PATH.
function hilbish.prependPath(dir) end

--- Changes the shell prompt to the provided string.
--- There are a few verbs that can be used in the prompt text.
--- These will be formatted and replaced with the appropriate values.
--- `%d` - Current working directory
--- `%u` - Name of current user
--- `%h` - Hostname of device
--- 
function hilbish.prompt(str, typ) end

--- Read input from the user, using Hilbish's line editor/input reader.
--- This is a separate instance from the one Hilbish actually uses.
--- Returns `input`, will be nil if Ctrl-D is pressed, or an error occurs.
function hilbish.read(prompt) end

--- Runs `cmd` in Hilbish's shell script interpreter.
--- The `streams` parameter specifies the output and input streams the command should use.
--- For example, to write command output to a sink.
--- As a table, the caller can directly specify the standard output, error, and input
--- streams of the command with the table keys `out`, `err`, and `input` respectively.
--- As a boolean, it specifies whether the command should use standard output or return its output streams.
--- 
function hilbish.run(cmd, streams) end

--- Sets the execution/runner mode for interactive Hilbish.
--- This determines whether Hilbish wll try to run input as Lua
--- and/or sh or only do one of either.
--- Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
--- sh, and lua. It also accepts a function, to which if it is passed one
--- will call it to execute user input instead.
--- Read [about runner mode](../features/runner-mode) for more information.
function hilbish.runnerMode(mode) end

--- Executed the `cb` function after a period of `time`.
--- This creates a Timer that starts ticking immediately.
function hilbish.timeout(cb, time) end

--- Checks if `name` is a valid command.
--- Will return the path of the binary, or a basename if it's a commander.
function hilbish.which(name) end

--- Puts a job in the background. This acts the same as initially running a job.
function hilbish.jobs:background() end

--- Puts a job in the foreground. This will cause it to run like it was
--- executed normally and wait for it to complete.
function hilbish.jobs:foreground() end

--- Evaluates `cmd` as Lua input. This is the same as using `dofile`
--- or `load`, but is appropriated for the runner interface.
function hilbish.runner.lua(cmd) end

--- Sets/toggles the option of automatically flushing output.
--- A call with no argument will toggle the value.
--- @param auto boolean|nil
function hilbish:autoFlush(auto) end

--- Flush writes all buffered input to the sink.
function hilbish:flush() end

--- Reads a liine of input from the sink.
--- @returns string
function hilbish:read() end

--- Reads all input from the sink.
--- @returns string
function hilbish:readAll() end

--- Writes data to a sink.
function hilbish:write(str) end

--- Writes data to a sink with a newline at the end.
function hilbish:writeln(str) end

--- Starts running the job.
function hilbish.jobs:start() end

--- Stops the job from running.
function hilbish.jobs:stop() end

--- Loads a module at the designated `path`.
--- It will throw if any error occurs.
function hilbish.module.load(path) end

--- Runs a command in Hilbish's shell script interpreter.
--- This is the equivalent of using `source`.
function hilbish.runner.sh(cmd) end

--- Starts a timer.
function hilbish.timers:start() end

--- Stops a timer.
function hilbish.timers:stop() end

--- Removes an alias.
function hilbish.aliases.delete(name) end

--- Get a table of all aliases, with string keys as the alias and the value as the command.
--- 
--- 
function hilbish.aliases.list() end

--- Resolves an alias to its original command. Will thrown an error if the alias doesn't exist.
function hilbish.aliases.resolve(alias) end

--- Creates a new job. This function does not run the job. This function is intended to be
--- used by runners, but can also be used to create jobs via Lua. Commanders cannot be ran as jobs.
--- 
--- 
function hilbish.jobs.add(cmdstr, args, execPath) end

--- Returns a table of all job objects.
function hilbish.jobs.all() end

--- Disowns a job. This simply deletes it from the list of jobs without stopping it.
function hilbish.jobs.disown(id) end

--- Get a job object via its ID.
--- @param id number
--- @returns Job
function hilbish.jobs.get(id) end

--- Returns the last added job to the table.
function hilbish.jobs.last() end

--- Adds a command to the history.
function hilbish.history.add(cmd) end

--- Retrieves all history as a table.
function hilbish.history.all() end

--- Deletes all commands from the history.
function hilbish.history.clear() end

--- Retrieves a command from the history based on the `index`.
function hilbish.history.get(index) end

--- Returns the amount of commands in the history.
function hilbish.history.size() end

--- Creates a timer that runs based on the specified `time`.
function hilbish.timers.create(type, time, callback) end

--- Retrieves a timer via its ID.
function hilbish.timers.get(id) end

return hilbish

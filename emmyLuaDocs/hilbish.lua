--- @meta

local hilbish = {}

--- Sets an alias of `cmd` to `orig`
--- @param cmd string
--- @param orig string
function hilbish.alias(cmd, orig) end

--- Appends `dir` to $PATH
--- @param dir string|table
function hilbish.appendPath(dir) end

--- Registers a completion handler for `scope`.
--- A `scope` is currently only expected to be `command.<cmd>`,
--- replacing <cmd> with the name of the command (for example `command.git`).
--- `cb` must be a function that returns a table of "completion groups."
--- Check `doc completions` for more information.
--- @param scope string
--- @param cb function
function hilbish.complete(scope, cb) end

--- Returns the current directory of the shell
function hilbish.cwd() end

--- Replaces running hilbish with `cmd`
--- @param cmd string
function hilbish.exec(cmd) end

--- Puts `fn` in a goroutine
--- @param fn function
function hilbish.goro(fn) end

--- Line highlighter handler. This is mainly for syntax highlighting, but in
--- reality could set the input of the prompt to *display* anything. The
--- callback is passed the current line and is expected to return a line that
--- will be used as the input display.
--- @param line string
function hilbish.highlighter(line) end

--- The command line hint handler. It gets called on every key insert to
--- determine what text to use as an inline hint. It is passed the current
--- line and cursor position. It is expected to return a string which is used
--- as the text for the hint. This is by default a shim. To set hints,
--- override this function with your custom handler.
--- @param line string
--- @param pos int
function hilbish.hinter(line, pos) end

--- Sets the input mode for Hilbish's line reader. Accepts either emacs or vim
--- @param mode string
function hilbish.inputMode(mode) end

--- Runs the `cb` function every `time` milliseconds.
--- Returns a `timer` object (see `doc timers`).
--- @param cb function
--- @param time number
--- @return table
function hilbish.interval(cb, time) end

--- Changes the continued line prompt to `str`
--- @param str string
function hilbish.multiprompt(str) end

--- Prepends `dir` to $PATH
--- @param dir string
function hilbish.prependPath(dir) end

--- Changes the shell prompt to `str`
--- There are a few verbs that can be used in the prompt text.
--- These will be formatted and replaced with the appropriate values.
--- `%d` - Current working directory
--- `%u` - Name of current user
--- `%h` - Hostname of device
--- @param str string
--- @param typ string Type of prompt, being left or right. Left by default.
function hilbish.prompt(str, typ) end

--- Read input from the user, using Hilbish's line editor/input reader.
--- This is a separate instance from the one Hilbish actually uses.
--- Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
--- @param prompt string
function hilbish.read(prompt) end

--- Runs `cmd` in Hilbish's sh interpreter.
--- If returnOut is true, the outputs of `cmd` will be returned as the 2nd and
--- 3rd values instead of being outputted to the terminal.
--- @param cmd string
function hilbish.run(cmd) end

--- Sets the execution/runner mode for interactive Hilbish. This determines whether
--- Hilbish wll try to run input as Lua and/or sh or only do one of either.
--- Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
--- sh, and lua. It also accepts a function, to which if it is passed one
--- will call it to execute user input instead.
--- @param mode string|function
function hilbish.runnerMode(mode) end

--- Runs the `cb` function after `time` in milliseconds
--- Returns a `timer` object (see `doc timers`).
--- @param cb function
--- @param time number
--- @return table
function hilbish.timeout(cb, time) end

--- Checks if `name` is a valid command
--- @param binName string
function hilbish.which(binName) end

return hilbish

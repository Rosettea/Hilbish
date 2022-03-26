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
--- A completion group is a table with the keys `items` and `type`.
--- `items` being a table of items and `type` being the display type of
--- `grid` (the normal file completion display) or `list` (with a description)
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

--- Sets the hinter function. This will be called on every key insert to determine
--- what text to use as an inline hint. The callback is passed 2 arguments:
--- the current line and the position. It is expected to return a string
--- which will be used for the hint.
--- @param cb function
function hilbish.hinter(cb) end

--- Sets the input mode for Hilbish's line reader. Accepts either emacs for vim
--- @param mode string
function hilbish.inputMode(mode) end

--- Runs the `cb` function every `time` milliseconds
--- @param cb function
--- @param time number
function hilbish.interval(cb, time) end

--- Changes the continued line prompt to `str`
--- @param str string
function hilbish.mlprompt(str) end

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
function hilbish.prompt(str) end

--- Read input from the user, using Hilbish's line editor/input reader.
--- This is a separate instance from the one Hilbish actually uses.
--- Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
--- @param prompt string
function hilbish.read(prompt) end

--- Runs `cmd` in Hilbish's sh interpreter.
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
--- @param cb function
--- @param time number
function hilbish.timeout(cb, time) end

--- Searches for an executable called `binName` in the directories of $PATH
--- @param binName string
function hilbish.which(binName) end

return hilbish

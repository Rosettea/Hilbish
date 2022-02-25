--- @meta

local hilbish = {}

--- Sets an alias of `orig` to `cmd`
function hilbish.alias() end

--- Appends `dir` to $PATH
function hilbish.appendPath() end

--- Registers a completion handler for `scope`.
--- A `scope` is currently only expected to be `command.<cmd>`,
--- replacing <cmd> with the name of the command (for example `command.git`).
--- `cb` must be a function that returns a table of the entries to complete.
--- Nested tables will be used as sub-completions.
function hilbish.complete() end

--- Returns the current directory of the shell
function hilbish.cwd() end

--- Replaces running hilbish with `cmd`
function hilbish.exec() end

--- Checks if the `f` flag has been passed to Hilbish.
function hilbish.flag() end

--- Puts `fn` in a goroutine
function hilbish.goroutine() end

--- Runs the `cb` function every `time` milliseconds
function hilbish.interval() end

--- Changes the continued line prompt to `str`
function hilbish.mlprompt() end

--- Prepends `dir` to $PATH
function hilbish.prependPath() end

--- Changes the shell prompt to `str`
--- There are a few verbs that can be used in the prompt text.
--- These will be formatted and replaced with the appropriate values.
--- `%d` - Current working directory
--- `%u` - Name of current user
--- `%h` - Hostname of device
function hilbish.prompt() end

--- Read input from the user, using Hilbish's line editor/input reader.
--- This is a separate instance from the one Hilbish actually uses.
--- Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
function hilbish.read() end

--- Runs `cmd` in Hilbish's sh interpreter.
function hilbish.run() end

--- Runs the `cb` function after `time` in milliseconds
function hilbish.timeout() end

return hilbish

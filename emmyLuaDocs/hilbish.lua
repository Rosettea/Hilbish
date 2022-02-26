--- @meta

local hilbish = {}

--- Sets an alias of `orig` to `cmd`
--- @param cmd string
--- @param orig string
function hilbish.alias(cmd, orig) end

--- Appends `dir` to $PATH
--- @param dir string|table
function hilbish.appendPath(dir) end

--- Registers a completion handler for `scope`.
--- A `scope` is currently only expected to be `command.<cmd>`,
--- replacing <cmd> with the name of the command (for example `command.git`).
--- `cb` must be a function that returns a table of the entries to complete.
--- Nested tables will be used as sub-completions.
function hilbish.complete() end

--- Returns the current directory of the shell
function hilbish.cwd() end

--- Replaces running hilbish with `cmd`
--- @param cmd string
function hilbish.exec(cmd) end

--- Puts `fn` in a goroutine
--- @param fn function
function hilbish.goroutine(fn) end

--- Runs the `cb` function every `time` milliseconds
--- @param cb function
--- @param time number
function hilbish.interval(cb, time) end

--- Changes the continued line prompt to `str`
--- @param str string
function hilbish.mlprompt(str) end

--- Prepends `dir` to $PATH
function hilbish.prependPath() end

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

--- Runs the `cb` function after `time` in milliseconds
--- @param cb function
--- @param time number
function hilbish.timeout(cb, time) end

return hilbish

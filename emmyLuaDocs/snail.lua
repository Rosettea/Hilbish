--- @meta

local snail = {}

--- Changes the directory of the snail instance.
--- The interpreter keeps its set directory even when the Hilbish process changes
--- directory, so this should be called on the `hilbish.cd` hook.
function snail:dir(path) end

--- Creates a new Snail instance.
function snail.new() end

--- Runs a shell command. Works the same as `hilbish.run`, but only accepts a table of streams.
function snail:run(command, streams) end

return snail

--- @meta

local snail = {}

--- Creates a new Snail instance.
function snail.new() end

--- Runs a shell command. Works the same as `hilbish.run`.
function snail:run(command, streams) end

return snail

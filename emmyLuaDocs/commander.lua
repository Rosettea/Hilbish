--- @meta

local commander = {}

--- Deregisters any command registered with `name`
function commander.deregister() end

--- Register a command with `name` that runs `cb` when ran
function commander.register() end

return commander

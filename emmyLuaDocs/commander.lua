--- @meta

local commander = {}

--- Deregisters any command registered with `name`
--- @param name string
function commander.deregister(name) end

--- Register a command with `name` that runs `cb` when ran
--- @param name string
--- @param cb function
function commander.register(name, cb) end

return commander

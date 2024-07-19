--- @meta

local commander = {}

--- Removes the named command. Note that this will only remove Commander-registered commands.
function commander.deregister(name) end

--- Adds a new command with the given `name`. When Hilbish has to run a command with a name,
--- it will run the function providing the arguments and sinks.
--- 
--- 
function commander.register(name, cb) end

--- Returns all registered commanders. Returns a list of tables with the following keys:
--- - `exec`: The function used to run the commander. Commanders require args and sinks to be passed.
function commander.registry() end

return commander

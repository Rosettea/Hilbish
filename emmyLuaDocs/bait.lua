--- @meta

local bait = {}

--- Catches a hook with `name`. Runs the `cb` when it is thrown
function bait.catch(name, cb) end

--- Same as catch, but only runs the `cb` once and then removes the hook
--- @param name string
--- @param cb function
function bait.catchOnce(name, cb) end

--- Returns a table with hooks (callback functions) on the event with `name`.
--- @param name string
--- @returns table<function>
function bait.hooks(name) end

--- Removes the `catcher` for the event with `name`.
--- For this to work, `catcher` has to be the same function used to catch
--- an event, like one saved to a variable.
--- @param name string
--- @param catcher function
function bait.release(name, catcher) end

--- Throws a hook with `name` with the provided `args`
--- @param name string
--- @vararg any
function bait.throw(name, ...args) end

return bait

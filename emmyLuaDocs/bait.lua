--- @meta

local bait = {}

--- Catches a hook with `name`. Runs the `cb` when it is thrown
--- @param name string
--- @param cb function
function bait.catch(name, cb) end

--- Same as catch, but only runs the `cb` once and then removes the hook
--- @param name string
--- @param cb function
function bait.catchOnce(name, cb) end

--- Removes the `catcher` for the event with `name`
--- For this to work, `catcher` has to be the same function used to catch
--- an event, like one saved to a variable.
function bait.release() end

--- Throws a hook with `name` with the provided `args`
--- @param name string
--- @vararg any
function bait.throw(name, ...) end

return bait

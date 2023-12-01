--- @meta

local bait = {}

--- Catches an event. This function can be used to act on events.
--- 
--- 
function bait.catch(name, cb) end

--- Catches an event, but only once. This will remove the hook immediately after it runs for the first time.
function bait.catchOnce(name, cb) end

--- Returns a list of callbacks that are hooked on an event with the corresponding `name`.
function bait.hooks(name) end

--- Removes the `catcher` for the event with `name`.
--- For this to work, `catcher` has to be the same function used to catch
--- an event, like one saved to a variable.
--- 
--- 
function bait.release(name, catcher) end

--- Throws a hook with `name` with the provided `args`.
function bait.throw(name, ...args) end

return bait

--- @meta

local bait = {}

--- Catches a hook with `name`. Runs the `cb` when it is thrown
function bait.catch() end

--- Same as catch, but only runs the `cb` once and then removes the hook
function bait.catchOnce() end

--- Throws a hook with `name` with the provided `args`
function bait.throw() end

return bait

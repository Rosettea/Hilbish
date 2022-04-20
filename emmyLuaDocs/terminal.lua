--- @meta

local terminal = {}

--- Restores the last saved state of the terminal
function terminal.restoreState() end

--- Saves the current state of the terminal
function terminal.saveState() end

--- Puts the terminal in raw mode
function terminal.setRaw() end

--- Gets the dimensions of the terminal. Returns a table with `width` and `height`
--- Note: this is not the size in relation to the dimensions of the display
function terminal.size() end

return terminal

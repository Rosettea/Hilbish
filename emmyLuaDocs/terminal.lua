--- @meta

local terminal = {}

--- Restores the last saved state of the terminal
function terminal.restoreState() end

--- Saves the current state of the terminal.
function terminal.saveState() end

--- Puts the terminal into raw mode.
function terminal.setRaw() end

--- Gets the dimensions of the terminal. Returns a table with `width` and `height`
--- NOTE: The size refers to the amount of columns and rows of text that can fit in the terminal.
function terminal.size() end

return terminal

--- @meta

local readline = {}

--- Deletes characters in the line by the given amount.
function readline:deleteByAmount(amount) end

--- Returns the current input line.
function readline:getLine() end

--- Returns the text that is at the register.
function readline:getVimRegister(register) end

--- Inserts text into the Hilbish command line.
function readline:insert(text) end

--- Prints a message *before* the prompt without it being interrupted by user input.
function readline:log(text) end

--- Creates a new readline instance.
function readline.new() end

--- Reads input from the user.
function readline:read() end

--- Reads a keystroke from the user. This is in a format of something like Ctrl-L.
function readline:getChar() end

--- Sets the vim register at `register` to hold the passed text.
function readline:setVimRegister(register, text) end

return readline

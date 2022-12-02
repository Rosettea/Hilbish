--- @meta

local hilbish = {}

--- 
--- @param cmd string
--- @param orig string
function hilbish.alias(cmd, orig) end

--- 
--- @param dir string|table
function hilbish.appendPath(dir) end

--- 
--- @param scope string
--- @param cb function
function hilbish.complete(scope, cb) end

--- 
function hilbish.cwd() end

--- 
--- @param cmd string
function hilbish.exec(cmd) end

--- 
--- @param fn function
function hilbish.goro(fn) end

--- 
--- @param line string
function hilbish.highlighter(line) end

--- 
--- @param line string
--- @param pos int
function hilbish.hinter(line, pos) end

--- 
--- @param mode string
function hilbish.inputMode(mode) end

--- 
--- @param cb function
--- @param time number
--- @return table
function hilbish.interval(cb, time) end

--- 
--- @param str string
function hilbish.multiprompt(str) end

--- 
--- @param dir string
function hilbish.prependPath(dir) end

--- 
--- @param str string
--- @param typ string Type of prompt, being left or right. Left by default.
function hilbish.prompt(str, typ) end

--- 
--- @param prompt string
function hilbish.read(prompt) end

--- 
--- @param cmd string
function hilbish.run(cmd) end

--- 
--- @param mode string|function
function hilbish.runnerMode(mode) end

--- 
--- @param cb function
--- @param time number
--- @return table
function hilbish.timeout(cb, time) end

--- 
--- @param binName string
function hilbish.which(binName) end

return hilbish

--- @meta

local util = {}

--- 
function util.AbbrevHome changes the user's home directory in the path string to ~ (tilde) end

--- 
function util. end

--- 
function util.DoFile runs the contents of the file in the Lua runtime. end

--- 
function util.DoString runs the code string in the Lua runtime. end

--- directory.
function util.ExpandHome expands ~ (tilde) in the path, changing it to the user home end

--- 
function util. end

--- 
function util.ForEach loops through a Lua table. end

--- 
function util. end

--- a string and a closure.
function util.HandleStrCallback handles function parameters for Go functions which take end

--- 
function util. end

--- 
function util. end

--- 
function util.SetExports puts the Lua function exports in the table. end

--- It is accessible via the __docProp metatable. It is a table of the names of the fields.
function util.SetField sets a field in a table, adding docs for it. end

--- is one which has a metatable proxy to ensure no overrides happen to it.
--- It sets the field in the table and sets the __docProp metatable on the
--- user facing table.
function util.SetFieldProtected sets a field in a protected table. A protected table end

--- Sets/toggles the option of automatically flushing output.
--- A call with no argument will toggle the value.
--- @param auto boolean|nil
function util:autoFlush(auto) end

--- Flush writes all buffered input to the sink.
function util:flush() end

--- 
function util. end

--- Reads a liine of input from the sink.
--- @returns string
function util:read() end

--- Reads all input from the sink.
--- @returns string
function util:readAll() end

--- Writes data to a sink.
function util:write(str) end

--- Writes data to a sink with a newline at the end.
function util:writeln(str) end

--- 
function util. end

--- 
function util. end

--- 
function util. end

return util

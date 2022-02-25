--- @meta

local fs = {}

--- Changes directory to `dir`
--- @param dir string
function fs.cd(dir) end

--- Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
--- @param name string
--- @param recursive bool
function fs.mkdir(name, recursive) end

--- Returns a table of files in `dir`
--- @param dir string
--- @return table
function fs.readdir(dir) end

--- Returns info about `path`
--- @param path string
function fs.stat(path) end

return fs

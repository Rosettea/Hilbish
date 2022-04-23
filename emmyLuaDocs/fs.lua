--- @meta

local fs = {}

--- Gives an absolute version of `path`.
--- @param path string
function fs.abs(path) end

--- Changes directory to `dir`
--- @param dir string
function fs.cd(dir) end

--- Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
--- @param name string
--- @param recursive boolean
function fs.mkdir(name, recursive) end

--- Returns a table of files in `dir`
--- @param dir string
--- @return table
function fs.readdir(dir) end

--- Returns info about `path`
--- @param path string
function fs.stat(path) end

return fs

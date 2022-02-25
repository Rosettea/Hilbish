--- @meta

local fs = {}

--- Changes directory to `dir`
function fs.cd() end

--- Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
function fs.mkdir() end

--- Returns a table of files in `dir`
function fs.readdir() end

--- Returns info about `path`
function fs.stat() end

return fs

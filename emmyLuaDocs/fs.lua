--- @meta

local fs = {}

--- Gives an absolute version of `path`.
--- @param path string
function fs.abs(path) end

--- Gives the basename of `path`. For the rules,
--- see Go's filepath.Base
function fs.basename() end

--- Changes directory to `dir`
--- @param dir string
function fs.cd(dir) end

--- Returns the directory part of `path`. For the rules, see Go's
--- filepath.Dir
function fs.dir() end

--- Glob all files and directories that match the pattern.
--- For the rules, see Go's filepath.Glob
function fs.glob() end

--- Takes paths and joins them together with the OS's
--- directory separator (forward or backward slash).
function fs.join() end

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

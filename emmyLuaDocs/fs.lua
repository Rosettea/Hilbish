--- @meta

local fs = {}

--- Gives an absolute version of `path`.
--- @param path string
function fs.abs(path) end

--- Gives the basename of `path`. For the rules,
--- see Go's filepath.Base
function fs.basename(path) end

--- Changes directory to `dir`
--- @param dir string
function fs.cd(dir) end

--- Returns the directory part of `path`. For the rules, see Go's
--- filepath.Dir
--- @param path string
function fs.dir(path) end

--- Glob all files and directories that match the pattern.
--- For the rules, see Go's filepath.Glob
--- @param pattern string
function fs.glob(pattern) end

--- Takes paths and joins them together with the OS's
--- directory separator (forward or backward slash).
--- @vararg any
function fs.join(...) end

--- Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
--- @param name string
--- @param recursive boolean
function fs.mkdir(name, recursive) end

--- Returns a table of files in `dir`.
--- @param dir string
--- @return table
function fs.readdir(dir) end

--- Returns a table of info about the `path`.
--- It contains the following keys:
--- name (string) - Name of the path
--- size (number) - Size of the path
--- mode (string) - Permission mode in an octal format string (with leading 0)
--- isDir (boolean) - If the path is a directory
--- @param path string
--- @returns table
function fs.stat(path) end

return fs

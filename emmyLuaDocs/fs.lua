--- @meta

local fs = {}

--- Returns an absolute version of the `path`.
--- This can be used to resolve short paths like `..` to `/home/user`.
function fs.abs(path) end

--- Returns the "basename," or the last part of the provided `path`. If path is empty,
--- `.` will be returned.
function fs.basename(path) end

--- Changes Hilbish's directory to `dir`.
function fs.cd(dir) end

--- Returns the directory part of `path`. If a file path like
--- `~/Documents/doc.txt` then this function will return `~/Documents`.
function fs.dir(path) end

--- Match all files based on the provided `pattern`.
--- For the syntax' refer to Go's filepath.Match function: https://pkg.go.dev/path/filepath#Match
--- 
--- 
function fs.glob(pattern) end

--- Takes any list of paths and joins them based on the operating system's path separator.
--- 
--- 
function fs.join(...path) end

--- Creates a new directory with the provided `name`.
--- With `recursive`, mkdir will create parent directories.
--- 
--- 
function fs.mkdir(name, recursive) end

--- Returns a pair of connected files, also known as a pipe.
--- The type returned is a Lua file, same as returned from `io` functions.
function fs.fpipe() end

--- Returns a list of all files and directories in the provided path.
function fs.readdir(path) end

--- Returns the information about a given `path`.
--- The returned table contains the following values:
--- name (string) - Name of the path
--- size (number) - Size of the path in bytes
--- mode (string) - Unix permission mode in an octal format string (with leading 0)
--- isDir (boolean) - If the path is a directory
--- 
--- 
function fs.stat(path) end

return fs

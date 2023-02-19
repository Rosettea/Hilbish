+ `command.exit` -> code, cmdStr > Thrown when a command exits.
`code` is the exit code of the command, and `cmdStr` is the command that was run.

+ `command.not-found` -> cmdStr > Thrown when a command is not found.

+ `command.not-executable` -> cmdStr > Thrown when Hilbish attempts to run a file
that is not executable.

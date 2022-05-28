local currentRunner = 'hybrid'
local runnerHandler = {}
local runners = {}

--- Get a runner by name.
--- @param name string
--- @return table
function runnerHandler.get(name)
	local r = runners[name]

	if not r then
		error(string.format('runner %s does not exist', name))
	end

	return r
end

--- Adds a runner to the table of available runners. If runner is a table,
--- it must have the run function in it.
--- @param name string
--- @param runner function | table
function runnerHandler.add(name, runner)
	if type(name) ~= 'string' then
		error 'expected runner name to be a table'
	end

	if type(runner) == 'function' then
		runner = {run = runner} -- this probably looks confusing
	end

	if type(runner) ~= 'table' then
		error 'expected runner to be a table or function'
	end

	if runners[name] then
		error(string.format('runner %s already exists', name))
	end

	runnerHandler.set(name, runner)
end

--- Sets a runner by name. The runner table must have the run function in it.
--- @param name string
--- @param runner table
function runnerHandler.set(name, runner)
	if not runner.run or type(runner.run) ~= 'function' then
		error 'run function in runner missing'
	end

	runners[name] = runner
end

--- Executes cmd with a runner. If runnerName isn't passed, it uses
--- the user's current runner.
--- @param cmd string
--- @param runnerName string?
--- @return string, number, string
function runnerHandler.exec(cmd, runnerName)
	if not runnerName then runnerName = currentRunner end

	local r = runnerHandler.get(runnerName)

	return r.run(cmd)
end

-- lsp shut up
hilbish = hilbish
--- Sets the current interactive/command line runner mode.
--- @param name string
function runnerHandler.setCurrent(name)
	local r = runnerHandler.get(name)
	currentRunner = name

	hilbish.runner.setMode(r.run)
end

-- add functions to hilbish.runner
for k, v in pairs(runnerHandler) do hilbish.runner[k] = v end

runnerHandler.add('hybrid', function(input)
	local cmdStr = hilbish.aliases.resolve(input)

	local _, _, err = hilbish.runner.lua(cmdStr)
	if not err then
		return input, 0, nil
	end

	return hilbish.runner.sh(input)
end)

runnerHandler.add('hybridRev', function(input)
	local _, _, err = hilbish.runner.sh(input)
	if not err then
		return input, 0, nil
	end

	local cmdStr = hilbish.aliases.resolve(input)
	return hilbish.runner.lua(cmdStr)
end)

runnerHandler.add('lua', function(input)
	local cmdStr = hilbish.aliases.resolve(input)
	return hilbish.runner.lua(cmdStr)
end)

runnerHandler.add('sh', function(input)
	return hilbish.runner.sh(input)
end)


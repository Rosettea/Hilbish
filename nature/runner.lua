local currentRunner = 'hybrid'
local runnerHandler = {}
local runners = {}

function runnerHandler.get(name)
	local r = runners[name]

	if not r then
		error(string.format('runner %s does not exist', name))
	end

	return r
end

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

function runnerHandler.set(name, runner)
	if not runner.run or type(runner.run) ~= 'function' then
		error 'run function in runner missing'
	end

	runners[name] = runner
end

function runnerHandler.exec(cmd, runnerName)
	if not runnerName then runnerName = currentRunner end

	local r = runnerHandler.get(runnerName)

	return r.run(cmd)
end

-- lsp shut up
hilbish = hilbish
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


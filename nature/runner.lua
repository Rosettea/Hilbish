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


function runnerHandler.setCurrent(name)
	local defaultRunners = {
		'hybrid',
		'hybridRev',
		'lua',
		'sh'
	}
	if defaultRunners[name] then
		hilbish.runner.setMode(name)
		return
	end

	local r = runnerHandler.get(name)
	currentRunner = name

	hilbish.runner.setMode(r.run)
end

-- add functions to hilbish.runner
for k, v in ipairs(runnerHandler) do hilbish.runner[k] = v end


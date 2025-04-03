-- @module hilbish

--- Sets the execution/runner mode for interactive Hilbish.
--- **NOTE: This function is deprecated and will be removed in 3.0**
--- Use `hilbish.runner.setCurrent` instead.
--- This determines whether Hilbish wll try to run input as Lua
--- and/or sh or only do one of either.
--- Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
--- sh, and lua. It also accepts a function, to which if it is passed one
--- will call it to execute user input instead.
--- Read [about runner mode](../features/runner-mode) for more information.
-- @param mode string|function
function hilbish.runnerMode(mode)
	if type(mode) == 'string' then
		hilbish.runner.setCurrent(mode)
	elseif type(mode) == 'function' then
		hilbish.runner.set('_', {
			run = mode
		})
		hilbish.runner.setCurrent '_'
	else
		error('expected runner mode type to be either string or function, got', type(mode))
	end
end

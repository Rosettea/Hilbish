---
title: Module hilbish.timers
description: timeout and interval API
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

If you ever want to run a piece of code on a timed interval, or want to wait
a few seconds, you don't have to rely on timing tricks, as Hilbish has a
timer API to set intervals and timeouts.

These are the simple functions `hilbish.interval` and `hilbish.timeout` (doc
accessible with `doc hilbish`, or `Module hilbish` on the Website).

An example of usage:
```lua
local t = hilbish.timers.create(hilbish.timers.TIMEOUT, 5000, function()
	print 'hello!'
end)

t:start()
print(t.running) // true
```

## Functions
|||
|----|----|
|<a href="#timers.create">create(type, time, callback) -> @Timer</a>|Creates a timer that runs based on the specified `time`.|
|<a href="#timers.get">get(id) -> @Timer</a>|Retrieves a timer via its ID.|

## Static module fields
|||
|----|----|
|INTERVAL|Constant for an interval timer type|
|TIMEOUT|Constant for a timeout timer type|

<hr><div id='timers.create'>
<h4 class='heading'>
hilbish.timers.create(type, time, callback) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
<a href="#timers.create" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a timer that runs based on the specified `time`.  
#### Parameters
`number` **`type`**  
What kind of timer to create, can either be `hilbish.timers.INTERVAL` or `hilbish.timers.TIMEOUT`

`number` **`time`**  
The amount of time the function should run in milliseconds.

`function` **`callback`**  
The function to run for the timer.

</div>

<hr><div id='timers.get'>
<h4 class='heading'>
hilbish.timers.get(id) -> <a href="/Hilbish/docs/api/hilbish/hilbish.timers/#timer" style="text-decoration: none;" id="lol">Timer</a>
<a href="#timers.get" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Retrieves a timer via its ID.  
#### Parameters
`number` **`id`**  


</div>

## Types
<hr>

## Timer
The Job type describes a Hilbish timer.
## Object properties
|||
|----|----|
|type|What type of timer it is|
|running|If the timer is running|
|duration|The duration in milliseconds that the timer will run|


### Methods
#### start()
Starts a timer.

#### stop()
Stops a timer.


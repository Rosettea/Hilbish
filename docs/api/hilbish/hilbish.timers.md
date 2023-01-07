---
title: Interface hilbish.timers
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
accessible with `doc hilbish`). But if you want slightly more control over
them, there is the `hilbish.timers` interface. It allows you to get
a timer via ID and control them.

All functions documented with the `Timer` type refer to a Timer object.

An example of usage:
```
local t = hilbish.timers.create(hilbish.timers.TIMEOUT, 5000, function()
	print 'hello!'
end)

t:start()
print(t.running) // true
```

## Interface fields
- `INTERVAL`: Constant for an interval timer type
- `TIMEOUT`: Constant for a timeout timer type

## Functions
### create(type, time, callback) -> <a href="#timer" style="text-decoration: none;">Timer</a>
Creates a timer that runs based on the specified `time` in milliseconds.
The `type` can either be `hilbish.timers.INTERVAL` or `hilbish.timers.TIMEOUT`

### get(id) -> <a href="#timer" style="text-decoration: none;">Timer</a>
Retrieves a timer via its ID.

## Types
## Timer
The Job type describes a Hilbish timer.
### Properties
- `type`: What type of timer it is
- `running`: If the timer is running
- `duration`: The duration in milliseconds that the timer will run

### Methods
#### start()
Starts a timer.

#### stop()
Stops a timer.


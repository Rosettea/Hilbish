---
name: Interface hilbish.timers
description: timeout and interval API
layout: apidoc
---

## Introduction
The timers interface si one to easily set timeouts and intervals
to run functions after a certain time or repeatedly without using
odd tricks.

## Object properties
- `type`: What type of timer it is
- `running`: If the timer is running
- `duration`: The duration in milliseconds that the timer will run

## Functions
### start()
Starts a timer.

### stop()
Stops a timer.

### create(type, time, callback)
Creates a timer that runs based on the specified `time` in milliseconds.
The `type` can either be interval (value of 0) or timeout (value of 1).

### get(id)
Retrieves a timer via its ID.


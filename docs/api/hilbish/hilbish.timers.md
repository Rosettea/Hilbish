---
title: Interface hilbish.timers
description: timeout and interval API
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction
The timers interface si one to easily set timeouts and intervals
to run functions after a certain time or repeatedly without using
odd tricks.

## Interface fields
- `INTERVAL`: Constant for an interval timer type
- `TIMEOUT`: Constant for a timeout timer type

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
The `type` can either be `hilbish.timers.INTERVAL` or `hilbish.timers.TIMEOUT`

### get(id) -> timer (Timer/Table)
Retrieves a timer via its ID.


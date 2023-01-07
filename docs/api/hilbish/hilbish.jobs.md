---
title: Interface hilbish.jobs
description: background job management
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

Manage interactive jobs in Hilbish via Lua.

Jobs are the name of background tasks/commands. A job can be started via
interactive usage or with the functions defined below for use in external runners.

## Types
### Job
Job Type.
#### Properties
- `cmd`: The user entered command string for the job.
- `running`: Whether the job is running or not.
- `id`: The ID of the job in the job table
- `pid`: The Process ID
- `exitCode`: The last exit code of the job.
- `stdout`: The standard output of the job. This just means the normal logs of the process.
- `stderr`: The standard error stream of the process. This (usually) includes error messages of the job.

## Functions
### background()
Puts a job in the background. This acts the same as initially running a job.

### foreground()
Puts a job in the foreground. This will cause it to run like it was
executed normally and wait for it to complete.

### start()
Starts running the job.

### stop()
Stops the job from running.

### add(cmdstr, args, execPath)
Adds a new job to the job table. Note that this does not immediately run it.

### all() -> table\<<a href="#job" style="text-decoration: none;">Job</a>>
Returns a table of all job objects.

### disown(id)
Disowns a job. This deletes it from the job table.

### get(id) -> <a href="#job" style="text-decoration: none;">Job</a>
Get a job object via its ID.

### last() -> <a href="#job" style="text-decoration: none;">Job</a>
Returns the last added job from the table.


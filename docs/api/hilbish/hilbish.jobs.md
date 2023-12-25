---
title: Module hilbish.jobs
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

## Functions
|||
|----|----|
|<a href="#jobs.add">add(cmdstr, args, execPath)</a>|Creates a new job. This function does not run the job. This function is intended to be|
|<a href="#jobs.all">all() -> table[@Job]</a>|Returns a table of all job objects.|
|<a href="#jobs.disown">disown(id)</a>|Disowns a job. This simply deletes it from the list of jobs without stopping it.|
|<a href="#jobs.get">get(id) -> @Job</a>|Get a job object via its ID.|
|<a href="#jobs.last">last() -> @Job</a>|Returns the last added job to the table.|

<hr>
<div id='jobs.add'>
<h4 class='heading'>
hilbish.jobs.add(cmdstr, args, execPath)
<a href="#jobs.add" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Creates a new job. This function does not run the job. This function is intended to be  
used by runners, but can also be used to create jobs via Lua. Commanders cannot be ran as jobs.  

#### Parameters
`string` **`cmdstr`**  
String that a user would write for the job

`table` **`args`**  
Arguments for the commands. Has to include the name of the command.

`string` **`execPath`**  
Binary to use to run the command. Does not

#### Example
```lua
hilbish.jobs.add('go build', {'go', 'build'}, '/usr/bin/go')
```
</div>

<hr>
<div id='jobs.all'>
<h4 class='heading'>
hilbish.jobs.all() -> table[<a href="/Hilbish/docs/api/hilbish/hilbish.jobs/#job" style="text-decoration: none;" id="lol">Job</a>]
<a href="#jobs.all" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns a table of all job objects.  

#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='jobs.disown'>
<h4 class='heading'>
hilbish.jobs.disown(id)
<a href="#jobs.disown" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Disowns a job. This simply deletes it from the list of jobs without stopping it.  

#### Parameters
`number` **`id`**  


</div>

<hr>
<div id='jobs.get'>
<h4 class='heading'>
hilbish.jobs.get(id) -> <a href="/Hilbish/docs/api/hilbish/hilbish.jobs/#job" style="text-decoration: none;" id="lol">Job</a>
<a href="#jobs.get" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Get a job object via its ID.  

#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='jobs.last'>
<h4 class='heading'>
hilbish.jobs.last() -> <a href="/Hilbish/docs/api/hilbish/hilbish.jobs/#job" style="text-decoration: none;" id="lol">Job</a>
<a href="#jobs.last" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the last added job to the table.  

#### Parameters
This function has no parameters.  
</div>

## Types
<hr>

## Job
The Job type describes a Hilbish job.
## Object properties
|||
|----|----|
|cmd|The user entered command string for the job.|
|running|Whether the job is running or not.|
|id|The ID of the job in the job table|
|pid|The Process ID|
|exitCode|The last exit code of the job.|
|stdout|The standard output of the job. This just means the normal logs of the process.|
|stderr|The standard error stream of the process. This (usually) includes error messages of the job.|


### Methods
#### background()
Puts a job in the background. This acts the same as initially running a job.

#### foreground()
Puts a job in the foreground. This will cause it to run like it was
executed normally and wait for it to complete.

#### start()
Starts running the job.

#### stop()
Stops the job from running.


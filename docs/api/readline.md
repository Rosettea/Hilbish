---
title: Module readline
description: Package readline is a pure-Go re-imagining of the UNIX readline API
layout: doc
menu:
  docs:
    parent: "API"
---

## Introduction

This package is designed to be run independently from murex and at some
point it will be separated into it's own git repository (at a stage when I
am confident that murex will no longer be the primary driver for features,
bugs or other code changes)

line reader library
The readline module is responsible for reading input from the user.
The readline module is what Hilbish uses to read input from the user,
including all the interactive features of Hilbish like history search,
syntax highlighting, everything. The global Hilbish readline instance
is usable at `hilbish.editor`.

Package terminal provides support functions for dealing with terminals, as
commonly found on UNIX systems.

Putting a terminal into raw mode is the most common requirement:

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
	        panic(err)
	}
	defer terminal.Restore(0, oldState)

Package terminal provides support functions for dealing with terminals, as
commonly found on UNIX systems.

Putting a terminal into raw mode is the most common requirement:

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
	        panic(err)
	}
	defer terminal.Restore(0, oldState)


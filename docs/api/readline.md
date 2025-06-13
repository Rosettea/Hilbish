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


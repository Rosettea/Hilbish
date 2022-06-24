---
title: Getting Started
layout: doc
weight: -10
---
# Contents
- [Installing Hilbish](#installing-hilbish)
- [Installing prebuilt binaries](#installing-prebuilt-binaries)
- [Building Hilbish](#building-hilbish)
- [Getting started](#getting-started)
- [Ending](#ending)

# Installing Hilbish
To get started, you have to first install Hilbish in your system. To do that you can either download the prebuilt binaries or build it yourself.

# Installing prebuilt binaries
You can download the latest Hilbish prebuilt from the [Github repository](https://github.com/Rosettea/Hilbish/releases/latest).

# Building Hilbish
 Prerequisites: 
 - Go 1.17+

After you've installed Go 1.17 or higher, you then can clone the Hilbish repository by using  
`git clone --recursive https://github.com/Rosettea/Hilbish`  
`cd Hilbish`  
`go get -d ./...`

The ``--recursive`` argument is required for the Lua libraries packaged with Hilbish.  
After you've cloned the whole repository and downloaded all needed packages, you then can run  

``git checkout $(git describe --tags `git rev-list --tags --max-count=1`)``  
`make build`  

to install the stable release of Hilbish.
However, if you wish to install the development release, you need to run  

`make dev`  

instead.

After you've built Hilbish, run  

`sudo make install`  

to install it globally.

# Getting started
After you've installed Hilbish you can then run it with

`hilbish`  

When Hilbish has booted for the first time, a message will show up saying you should run the 

`guide`  

command. The guide in question is very simple, and should help you start your journey with Hilbish. Proper documentation can be found after running

`doc`

Hilbish is the same as any other Linux shell, with the addition of running Lua. It can also act as a Lua REPL by running   

`hilbish.runnerMode 'lua'`.  

The REPL mode can be turned off via  

`hilbish.runnerMode 'hybrid`.  

# Ending
Now you're ready to use Hilbish. Feel free to check the documentation and guide for more information. Good luck!

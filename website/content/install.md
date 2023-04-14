---
title: Install
description: Steps on how to install Hilbish on all the OSes and distros supported.
layout: page
---

## Official Binaries
The best way to get Hilbish is to get a build directly from GitHub.
At any time, there are 2 versions of Hilbish recommended for download:
the latest stable release, and development builds from the master branch.

You can download both at any time, but note that the development builds may
have breaking changes.

For the latest **stable release**, check here: https://github.com/Rosettea/Hilbish/releases/latest  
For a **development build**: https://nightly.link/Rosettea/Hilbish/workflows/build/master

## Compiling
To read the steps for compiling Hilbish, head over to the [GitHub repository.](https://github.com/Rosettea/Hilbish#build)

## Package Repositories
Methods of installing Hilbish for your Linux distro.

### Fedora (COPR)
An official COPR is offered to install Hilbish easily on Fedora.
Enable the repo:
```
sudo dnf copr enable sammyette/Hilbish
```

And install Hilbish:
```
sudo dnf install hilbish
```

Or for the latest development build from master:
```
sudo dnf install hilbish-git
```

### Arch Linux (AUR)
Hilbish is on the AUR. Setup an AUR helper, and install.
Example with yay:  

```
yay -S hilbish
```

Or, from master branch:  
```
yay -S hilbish-git
```

### Alpine Linux
Hilbish is currentlty in the testing/edge repository for Alpine.
Follow the steps [here](https://wiki.alpinelinux.org/wiki/Enable_Community_Repository)
(Using testing repositories) and install:  
```
apk add hilbish
```

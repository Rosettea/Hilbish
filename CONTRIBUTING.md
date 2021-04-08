# Contributing to Hilbish
Thanks for your interest in contributing to Hilbish! No matter if it's
a new feature, documentation, or bug report, we appreciate any kind
of contributions.

This file is a document to state the steps/rules to take when making
a contribution. Be sure to read through it.

## Bug Reports // Feature Requests
Use GitHub Issues to report any bugs or to request any features
that may be useful to *anyone* else.

Check [currently open issues](https://github.com/Hilbis/Hilbish/issues)
and [closed ones](https://github.com/Hilbis/Hilbish/issues?q=is%3Aissue+is%3Aclosed) to make sure someone else hasn't already made the issue.

For bug reports, be sure to include:
- Hilbish Version (`hilbish -v`)
- Ways to reproduce

## Code
For any code contributions (Lua and/or Go), you should follow these
rules:  
- Tab size 4 indentation
- In Lua prefer no braces `()` if the function takes 1 argument
- Use camelCase for function names

### Making the Pull Request
1. Ensure that any new install or build dependencies are documented in
the README.md and PR request.

2. Say in the pull request details the changes to the shell,
this includes useful file locations and breaking changes.

3. The versioning scheme we use is [SemVer](http://semver.org/) and the
commit scheme we use is
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
Please document any backwards incompatible changes and be sure to name
your commits correctly.

4. Finally, make the pull request to the **dev** branch.

## Finding Issues to Contribute to
You can check out the [help wanted](https://github.com/Hilbis/Hilbish/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22+)
labels to figure out what we need your help working on.  

The [up for grabs](https://github.com/Hilbis/Hilbish/issues?q=is%3Aissue+is%3Aopen+label%3A%22up+for+grabs%22+) labeled issues are low hanging fruit that should be
easy for anyone. You can use this to get started on contributing!

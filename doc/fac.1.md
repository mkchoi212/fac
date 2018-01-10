# NAME

fac - tool for resolving git merge conflicts

# SYNOPSIS

fac

# DESCRIPTION

Fix All Conflicts (fac) aims to make resolving git(1) merge conflicts easier.
It provides an ncurses-based based terminal user interface to resolve conflicts
interactively.

The UI is split into three panes. Two show the versions of content in conflict.
The third, a sidebar show actions that can be taken to resolve the conflict.

# USAGE

fac operates much like `git add -p`. It has a prompt for input at the bottom of
the screen where the various commands are entered.

The commands have been preset to the following specifications:

**w** - show more lines up  
**s** - show more lines down  
**a** - use local version  
**d** - use incoming version

**j** - scroll down  
**k** - scroll up

**v** - [v]iew orientation  
**n** - [n]ext conflict  
**p** - [p]revious conflict

**h** | **?** - [h]elp  
**q** | **Ctrl+c** - [q]uit

The movement controls have been derived from both the world of gamers (WASD)
and vi(1) users (HJKL).

# BUGS

fac does not currently support diff3 merge conflict output.

# NOTES

1. Home page and source code  
   https://github.com/mkchoi212/fac

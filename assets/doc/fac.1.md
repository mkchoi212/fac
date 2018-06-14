# NAME

fac - tool for resolving git merge conflicts

# SYNOPSIS

**fac -h | -help**\
**fac -version**\

# DESCRIPTION

Fix All Conflicts (fac) aims to make resolving git(1) merge conflicts easier.
It provides an ncurses-based based terminal user interface to resolve conflicts
interactively.

The UI is split into three panes. Two show the versions of content in conflict.
The third is a sidebar that show the queue and actions that can be taken to resolve the conflict.

# USAGE

fac operates much like `git add -p`. It has a prompt for input at the bottom of
the screen where the various commands are entered.

The commands have been preset to the following specifications:

**w** 
: show more lines up  

**s** 
: show more lines down

**a**
: use local version  

**d**
: use incoming version

**e**
: manually edit code

**j**
: scroll down  

**k**
: scroll up

**v** 
: [v]iew orientation  

**n**
: [n]ext conflict  

**p**
: [p]revious conflict

**h | ?**
: [h]elp  

**q | Ctrl+c**
: [q]uit

The movement controls have been derived from both the world of gamers (WASD)
and vi(1) users (HJKL). To customize these key bindings, please refer to the next section.

# CUSTOMIZATION

Program behavior can be customized by creating a `$HOME/.fac.yml` with the below variables.

## BEHAVIOR

**cont_eval**
: evaluate commands without pressing ENTER

## KEY BINDINGS

**select_local**
: select local version

**select_incoming**
: select incoming version

**toggle_view**
: toggle to horizontal | vertical view

**show_up**   
: show more lines above

**show_down**
: show more lines below

**scroll_up**
: scroll up

**scroll_down**
: scroll down

**edit**
: manually edit code chunk

**next**
: go to next conflict

**previous**
: go to previous conflict

**quit**
: quit application

**help**
: display help in side bar

# COPYRIGHT

Copyright (c) 2018 Mike JS. Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

# AUTHOR

Mike JS. Choi <mkchoi212@icloud.com>\
shogo-ma      <choroma194@gmail.com>

Please send bug reports or comments to <https://github.com/mkchoi212/fac/issues>.\
For more information, see the homepage at <https://github.com/mkchoi212/fac>.


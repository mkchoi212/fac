package main

import (
	"flag"
	"fmt"
	"os"
)

// ParseFlags parses flags provided by the user
func ParseFlags() {
	// Setup custom help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:

   fac

Customizable variables: 
  Behavior

	cont_eval        evaluate commands without pressing ENTER

  Key bindings

	select_local     select local version
	select_incoming  select incoming version
	toggle_view      toggle to horizontal | horizontal view
	show_up          show more lines above
	show_down        show more lines below
	scroll_up        ...
	scroll_down      ...
	edit             manually edit code chunk
	next             go to next conflict
	previous         go to previous conflict
	quit             ...
	help             display help in side bar

Following variables may be defined in your $HOME/.fac.yml to customize behavior

`)
	}

	flag.Parse()
}

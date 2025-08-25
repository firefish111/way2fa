package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/firefish111/way2fa/parse/csv"
	"github.com/firefish111/way2fa/ui"
	"os"
)

const (
	VersionMajor = 0
	VersionMinor = 2
	VersionPatch = 0
)

func main() {
	ver := flag.Bool("version", false, "Prints version")

	flag.Parse()

	// end of args

	if *ver {
		fmt.Printf("way2fa v%d.%d.%d\n", VersionMajor, VersionMinor, VersionPatch)
		return
	}

	var name *string
	if a := flag.Args(); len(a) > 0 {
		name = &a[0]
	}

	store, err := csv.GetFile(name)
	if err != nil {
		panic(err)
	}

	// call the ui
	model, err := ui.Create(*store)
	if err != nil {
		panic(err)
	}

	prog := tea.NewProgram(model)
	if _, err := prog.Run(); err != nil { // do the running
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

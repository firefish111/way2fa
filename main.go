package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/internal/ui"
)

const (
	VersionMajor = 0
	VersionMinor = 3
	VersionPatch = -1
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

	// deal with the fact that config dir may not be real
	err := config.InitConfPath()
	if err != nil {
		panic(err)
	}

	// store is the automatic detector
	store := detector.Detect(name)
	if store == nil {
		fmt.Fprintf(os.Stderr, "Automatic detection failed, aborting.\n\nHINT: create a file in %s.\n", config.ConfPath)
		return
	}

	// call the ui
	model, err := ui.Create(store)
	if err != nil {
		panic(err)
	}

	ctrller, err := bubblon.New(model)
	if err != nil {
		panic(err)
	}

	prog := tea.NewProgram(ctrller, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil { // do the running
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/internal/ui"
)

const (
	VersionMajor = 0
	VersionMinor = 3
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

	// deal with the fact that config dir may not be real
	err := config.InitConfPath()
	if err != nil {
		panic(err)
	}

	// store is the automatic detector
	store := detector.Detect(name)
	if store == nil {
		panic(fmt.Errorf("Automatic detection failed, aborting.\n\nHINT: create a file in %s.\n", config.ConfPath))
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
	doneModel, err := prog.Run() // as for why `doneModel`, see below

	if err != nil { // do the running
		panic(fmt.Errorf("Error running model:\n\t%w\n", err))
	}

	// tea.NewProgram copies the model it takes.
	// prog.Run() returns the model after it's done with it, but as an interface, which needs to be upcasted,
	// hence why we use another variable `doneModel`
	ctrller = doneModel.(bubblon.Controller)
	if ctrller.Err != nil { // when bubblon.Fail is called, the error is put in Err, so it can gracefully exit tea
		panic(fmt.Errorf("Runtime error during model execution:\n\t%w\n", ctrller.Err))
	}
}

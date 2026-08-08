package main

import (
	"flag"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/donderom/bubblon/v2"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/internal/ui"
	"github.com/firefish111/way2fa/manager"
)

const (
	VersionMajor = 0
	VersionMinor = 4
	VersionPatch = 1
)

func main() {
	ver := flag.Bool("version", false, "Prints version")
	export := flag.String("export", "", "Exports 2FA store to a different store, like so:\n\t-export dest [src]\nAlso necessary to change passwords.")
	create := flag.Bool("create", false, "Creates a new 2FA store, like so:\n-create [dest]")
	list := flag.Bool("list", false, "Lists the default files that are present")

	flag.Parse()

	// end of args

	if *ver {
		fmt.Printf("way2fa v%d.%d.%d\n", VersionMajor, VersionMinor, VersionPatch)
		return
	}

	if *list {
		manager.List()
		return
	}

	var name *string
	if a := flag.Args(); len(a) > 0 {
		name = &a[0]
	}

	if *export != "" {
		err := manager.Export(name, export)
		if err != nil {
			panic(fmt.Errorf("Export failed: %w", err))
		}
		return
	}

	if *create {
		_, err := manager.Create(name)
		if err != nil {
			panic(fmt.Errorf("Create failed: %w", err))
		}
		return
	}

	// deal with the fact that config dir may not be real
	err := config.InitConfPath()
	if err != nil {
		panic(err)
	}

	// store is the automatic detector
	store, err := detector.Detect(name)
	if err != nil {
		panic(fmt.Errorf("Automatic detection failed, aborting.\n\n%w\n\nHINT: provide a file, or use -create.\n", err))
	}

	// load accounts from file
	err = store.Load()
	if err != nil {
		panic(fmt.Errorf("Failed to load store: %w", err))
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

	prog := tea.NewProgram(ctrller)
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

	// save accounts to file.
	// done at last, after ui has exited
	err = store.Save()
	if err != nil {
		panic(fmt.Errorf("Failed to save store: %w", err))
	}
}

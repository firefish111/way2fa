package manager

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/ui/manager"
	"github.com/firefish111/way2fa/parse"
)

// Export an existing 2FA store to a new file.
// If either parameter is nil, that means the default file, which is source dependent.
// Most of the logic is in Create().
func Export(src, dest *string) error {
	if (src == nil && dest == nil) || (src != nil && dest != nil && *src == *dest) {
		return errors.New("Cannot export to and from the same file, or have both set to a default value.")
	}

	// detect source file
	srclist, err := detector.Detect(src)
	if err != nil {
		return err
	}

	// load accounts from source file
	err = srclist.Load()
	if err != nil {
		return err
	}

	// create new destination list
	destlist, err := Create(dest)
	if err != nil {
		return err
	}

	if srclist.IsPasswordProtected() {
		// TODO
		return errors.New("NOT IMPLEMENTED: exporting of a password-protected file.")
	}

	// TODO: srclist.Decrypt()
	// get accounts
	srcAccs, err := srclist.GetAccs()
	if err != nil {
		return err
	}

	// TODO: destlist.Decrypt()
	// store them in new destination
	err = destlist.SetAccs(srcAccs)
	if err != nil {
		return err
	}

	// save them to disk
	err = destlist.Save()
	if err != nil {
		return err
	}

	return nil
}

// Create a new 2FA store with the given name.
func Create(name *string) (parse.AccountList, error) {
	var reader parse.AccountList

	ctrller, err := bubblon.New(manager.CreateCreatorModel(name, &reader))
	if err != nil {
		panic(err)
	}

	prog := tea.NewProgram(ctrller)
	doneModel, err := prog.Run()

	if err != nil { // do the running
		return nil, fmt.Errorf("Error running model:\n\t%w\n", err)
	}

	// tea.NewProgram copies the model it takes.
	// prog.Run() returns the model after it's done with it, but as an interface, which needs to be upcasted,
	// hence why we use another variable `doneModel`
	ctrller = doneModel.(bubblon.Controller)
	if ctrller.Err != nil { // when bubblon.Fail is called, the error is put in Err, so it can gracefully exit tea
		return nil, fmt.Errorf("Runtime error during model execution:\n\t%w\n", ctrller.Err)
	}

	return reader, nil
}

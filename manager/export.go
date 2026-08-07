package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/internal/ui/password"
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

	// load acctouns from source file
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
	// FIXME: null deref of destlist
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
	var allTypes []parse.AccountList
	{
		ways := detector.GetPossibleWayFormats()
		pures := detector.GetPossiblePureFormats()

		allTypes = make([]parse.AccountList, 0, len(ways)+len(pures))
		for _, way := range ways {
			allTypes = append(allTypes, way)
		}
		for _, pure := range pures {
			allTypes = append(allTypes, pure)
		}
	}

	// TODO: replace by a proper bubbletea ui
	fmt.Println("Select destination store type: ")
	for i, typ := range allTypes { // print all of them
		fmt.Printf("\t[%d]: %s\n", i+1, reflect.TypeOf(typ).Elem().Name())
	}

	// get index
	var index int
	fmt.Print("> ")
	fmt.Scanf("%d", &index)

	// get reader
	reader := allTypes[index-1]

	// get default filename
	var filename string
	var err error

	if name != nil && *name != "" {
		filename, err = filepath.Abs(*name)
	} else if pure, ok := reader.(parse.PureAccountList); ok {
		filename = filepath.Join(
			config.ConfPath,
			pure.GetDefaultFilename(),
		)
	} else if _, ok := reader.(parse.WayAccountList); ok {
		filename = filepath.Join(
			config.ConfPath,
			detector.DefaultWayFilename,
		)
	} else {
		return nil, fmt.Errorf("Can't match account list #%d to filename %v", index, name)
	}

	// another if statement as the one before this is to test the defaults
	if way, ok := reader.(parse.WayAccountList); ok {
		var input rune
		fmt.Println("Set a password? (Y/n)")
		fmt.Scanf("%c", &input)

		// set it to password protected depending on input: Y is default so we only compare N
		way.SetPasswordProtected(!(input == 'N' || input == 'n'))

	}

	// exclusive create to avoid overwriting existing files
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("File %s already exists: %w", filename, err)
		}

		return nil, err
	}
	f.Close() // leave it to the interface itself to do that

	if reader.IsPasswordProtected() {
		// prepopulate from file
		reader.PrepopulateFromFile(filename, name == nil || *name == "")
		reader.PopulateNew()

		// tea model for password prompt - very similar to code in /main.go
		model := password.CreatePasswordPrompt(reader, "Set a password")
		ctrller, err := bubblon.New(model)
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
	}

	reader.Save()

	return reader, nil
}

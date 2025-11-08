// Package to detect the type of file.
//
// This uses the modules in parse: each subpackage of parse/modules is its own filetype.
package detector

import (
	"os"

	"github.com/firefish111/way2fa/parse"
	"github.com/firefish111/way2fa/parse/modules/csv_pure"
)

// This function gets a list of AccountLists with empty fields.
// These are designed to have the PrePopulate methods called on them.
//
// These are also in the priority order that they should be searched in.
// As of yet, this is just CsvPure at highest priority. XXX please update when necessary
func getAllUnpopulated() []parse.AccountList {
	return []parse.AccountList{
		&csv_pure.CsvPure{},
	}
}

func Detect(path *string) parse.AccountList {
	types := getAllUnpopulated()
	var err error

	for i, _ := range types {
		// if this errors, it usually means directory not found
		if path == nil {
			err = types[i].PrepopulateDefault()
		} else {
			err = types[i].PrepopulateWith(*path)
		}

		// henceforth, types[i] is no longer empty, unless err is nil

		if err != nil { // i.e. we have no chance at being existing
			continue
		}

		if _, err = os.Stat(types[i].GetSourceFilePath()); err != nil { // does not exist
			continue
		}
		
		if !types[i].Validate() { // if is invalid
			continue
		}

		return types[i] // found it!
	}

	// all failed, we got nothing
	return nil
}

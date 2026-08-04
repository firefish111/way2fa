package detector

import (
	"path/filepath"

	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/parse"
	csv_pure "github.com/firefish111/way2fa/parse/modules/pure/csv"
)

// This function gets a list of AccountLists with empty fields.
// These are designed to have the Prepopulate methods called on them.
//
// These are also in the priority order that they should be searched in.
// As of yet, this is just CsvPure at highest priority.
// XXX please update when necessary
func getPossiblePureFormats() []parse.PureAccountList {
	return []parse.PureAccountList{
		&csv_pure.CsvPure{},
	}
}

// behaves a little bit differently to .way format, as we have to ask each format individually
// what filename it's looking for.
func tryDetectPure(path_optional *string) parse.PureAccountList {
	possibilities := getPossiblePureFormats()

	for i, _ := range possibilities {
		// get filename we're looking for
		var path string
		if path_optional == nil {
			path = filepath.Join(
				config.ConfPath,
				possibilities[i].GetDefaultFilename(),
			)
		} else {
			path = *path_optional
		}

		err := possibilities[i].PrepopulateFromFile(path, path_optional == nil)
		if err != nil { // we don't return, as failing here is not an error, just a signal to move on to next one
			continue
		}

		// is can be detected
		if isValid := possibilities[i].Detect(); !isValid {
			continue
		}

		// we found it!
		return possibilities[i]
	}

	// after every loop, we found nothing
	return nil
}

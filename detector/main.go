// Package to detect the type of file.
//
// This uses the modules in parse: each subpackage of parse/modules is its own filetype.
package detector

import (
	"fmt"
	"os"

	"github.com/firefish111/way2fa/format"
	"github.com/firefish111/way2fa/parse"
)

// This takes a file as input, to verify the validity of.
// This is often nil if no args are given, in which case each individual file type uses its default filename
func Detect(path *string) (parse.AccountList, error) {
	// TODO: this is perhaps bad, as the `.way` header tells you what filetype it is.
	// TODO: perhaps add an IsPure() method and/or a GetHeaderIndicator() method to each AccountList instance
	// this would make checks faster, as all the checks are done here

	if path != nil {
		if _, err := os.Stat(*path); err != nil { // does not exist
			return nil, format.EmptyHandedError(fmt.Sprintf("provided path \"%s\" doesn't exist", *path))
		}
	}

	// highest priority is finding the (more secure) .way files.
	// only if we fail do we switch the plaintext pure files.
	way, err := tryDetectWay(path)
	if err == nil {
		// we found one!
		return way, nil
	}

	// if we got here, we failed to detect a .way file, so we fall back to the old pure file format.
	pure := tryDetectPure(path)
	if pure != nil {
		return pure, nil
	}

	if path != nil {
		return nil, format.EmptyHandedError(fmt.Sprintf("don't know what to do with provided file \"%s\"", *path))
	} else {
		return nil, format.EmptyHandedError("cannot find defaults")
	}
}

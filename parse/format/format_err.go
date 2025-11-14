package format

import (
	"fmt"
)

// Custom error type for problems with the .way format.
//
// Not to be confused with detection issues, those are non-erroring, and arise only in the detector.
type FormatError struct {
	reason string
}

func (e FormatError) Error() string {
	return fmt.Sprintf("FormatError: epic fail (%s)", e.reason)
}

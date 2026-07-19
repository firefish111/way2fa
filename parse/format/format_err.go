package format

import (
	"fmt"
)

// Custom error type for problems with the .way format.
//
// Not to be confused with detection issues, those are non-erroring, and arise only in the detector.
type FormatError struct {
	reason          string
	isFaultOfHeader bool
}

func (e FormatError) Error() string {
	if e.isFaultOfHeader {
		return fmt.Sprintf("FormatError: malformed account file header; %s", e.reason)
	} else {
		return fmt.Sprintf("FormatError: malformed account file; %s", e.reason)
	}
}

func headerError(reason string) FormatError {
	return FormatError{reason: reason, isFaultOfHeader: true}
}

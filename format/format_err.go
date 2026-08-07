package format

import (
	"fmt"
)

type (
	formatStatus int
	// Custom error type for problems with the .way format.
	//
	// Not to be confused with detection issues, those are non-erroring, and arise only in the detector.
	FormatError struct {
		reason string
		status formatStatus
	}
)

const (
	malformedHeader formatStatus = iota
	failedDetection
	emptyHanded
)

func (e FormatError) Error() string {
	switch e.status {
	case malformedHeader:
		return fmt.Sprintf("FormatError: malformed account file header; %s", e.reason)
	case failedDetection:
		return fmt.Sprintf("FormatError: detector failed; %s", e.reason)
	case emptyHanded:
		return fmt.Sprintf("Error: Found no valid file to read from; %s", e.reason)
	default:
		return fmt.Sprintf("FormatError: %s", e.reason)
	}
}

func headerError(reason string) FormatError {
	return FormatError{reason: reason, status: malformedHeader}
}

func DetectionError(reason string) FormatError {
	return FormatError{reason: reason, status: failedDetection}
}

func EmptyHandedError(reason string) FormatError {
	return FormatError{reason: reason, status: emptyHanded}
}

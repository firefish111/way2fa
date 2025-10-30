package csv_pure

import (
	"path/filepath"
)

// Whether this is a valid CSV constitutes a valid csv file.
// As to whether the csv file is actually valid is a "content error", not a "format issue".
//
// Hence, we naively check whether the extension is csv as to whether this is valid.
// This is horrendous, but nobody should be using this module anyway as it's horribly insecure, so cry about it
func (c CsvPure) Validate() bool {
	is_csv := filepath.Ext(c.path) == CsvPureExt
	return is_csv
}

package account

// Implement gocsv's CSVUnmarshaler
func (r *RawKey) UnmarshalCSV(value string) (err error) {
	*r, err = decodeTextKey(value)
	return
}

package account

// Implement rowboat's CSVUnmarshaler
func (r *RawKey) UnmarshalCSV(value string) (err error) {
	*r, err = decodeTextKey(value)
	return
}

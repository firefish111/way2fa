package account

// Implement gocsv's CSVUnmarshaler
func (r *RawKey) UnmarshalCSV(value string) (err error) {
	*r, err = DecodeTextKey(value)
	return
}

// Implement gocsv's CSVMarshaler
func (r *RawKey) MarshalCSV() (string, error) {
	out := EncodeTextKey(*r)
	return out, nil
}

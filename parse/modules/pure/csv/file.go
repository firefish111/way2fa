package csv_pure

import (
	"fmt"
	"io"
	"os"
)

func (c *CsvPure) Load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return fmt.Errorf("cannot access keyfile %s: %w", c.path, err)
	}

	defer f.Close() // wait until end of function to close

	buf, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("cannot read keyfile %s: %w", c.path, err)
	}

	c.buffer = string(buf)
	return nil
}

func (c *CsvPure) Save() error {
	// os.Open only opens readonly
	f, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("cannot create keyfile %s: %w", c.path, err)
	}
	defer f.Close()

	_, err = f.WriteString(c.buffer)
	if err != nil {
		return fmt.Errorf("cannot save keyfile %s: %w", c.path, err)
	}

	return nil
}

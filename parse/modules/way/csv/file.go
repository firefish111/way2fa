package csv_way

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (c *CsvWay) Load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// load header
	err = binary.Read(f, binary.LittleEndian, &c.header)
	if err != nil {
		return err
	}

	// validate header. should be valid thanks to detector, but can't be too sure
	if err = c.header.Validate(); err != nil {
		return err
	}

	// understand the header
	c.processHeader()

	// load cryptor
	err = binary.Read(f, binary.LittleEndian, &c.crypt)
	if err != nil {
		return err
	}

	// read payload
	c.payload = make([]byte, c.header.PayloadSize)
	n, err := io.ReadFull(f, c.payload)
	if err != nil {
		return fmt.Errorf("failed to read payload, could only read %d of expected %d bytes: %w", n, c.header.PayloadSize, err)
	}

	return nil
}

// opposite of Load()
func (c *CsvWay) Save() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// ensure header is meaningful
	c.updateHeader()

	// write header
	err = binary.Write(f, binary.LittleEndian, c.header)
	if err != nil {
		return err
	}

	// write cryptor
	err = binary.Write(f, binary.LittleEndian, c.crypt)
	if err != nil {
		return err
	}

	// write payload
	_, err = f.Write(c.payload)
	if err != nil {
		return err
	}

	return nil
}

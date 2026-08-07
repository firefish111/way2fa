package csv_way

import (
	"time"

	"github.com/firefish111/way2fa/format"
)

func (c *CsvWay) processHeader() {
	// interpret flags
	c.hasPassword = c.header.Flags&format.HasPassword != 0

	// unpack capabilities
	c.capabilities = c.header.Capabilities.Unpack()
}

func (c *CsvWay) updateHeader() {
	// a lot of these are redundant if we assume that header is already valid, but this may not be a guarantee

	// multi-layer cast (ugly)
	c.header.Magic = [4]byte([]byte(format.MagicNumber))

	// set latest version: we might have read from an older version
	c.header.Version = format.FormatVersion

	// set file type
	c.header.FileType = format.Csv

	// set flags
	c.header.Flags = 0
	if c.hasPassword {
		c.header.Flags |= format.HasPassword
	}

	// set last-modified time
	c.header.LastModified = time.Now().Unix()

	// write payload length
	c.header.PayloadSize = uint64(len(c.payload))

	// repack capabilities
	c.header.Capabilities = c.capabilities.Pack()
}

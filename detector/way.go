package detector

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/firefish111/way2fa/format"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/parse"
	csv_way "github.com/firefish111/way2fa/parse/modules/way/csv"
)

// all .way files have the same name by default for consistency's sake.
// exactly what format it uses is dictated by the file itself
const (
	defaultWayExt      = "way"
	defaultWayFilename = "keys" + defaultWayExt
)

// This function gets a list of AccountLists with empty fields.
// These are designed to have the Prepopulate methods called on them.
//
// XXX please update when necessary
func getPossibleWayFormats() map[format.FileTypeId]parse.WayAccountList {
	return map[format.FileTypeId]parse.WayAccountList{
		// TODO
		format.Csv: &csv_way.CsvWay{},
	}
}

// try to detect .way file, at the given location, or if nil, at the default location.
// returns an error saying whether it succeeded or not.
func tryDetectWay(path_optional *string) (parse.WayAccountList, error) {
	var path string
	if path_optional == nil {
		path = filepath.Join(
			config.ConfPath,
			defaultWayFilename,
		)
	} else {
		path = *path_optional
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// try to read the way file
	f, err := os.Open(path)
	if err != nil {
		return nil, err // file doesn't exist, we have failed
	}

	// NOTE: whereas ordinarily we would do this, we explicitly close later on so that file descriptor is not open for Prepopulate call
	// defer f.Close()

	header := format.Header{}
	// NOTE: we read this little endian, as format requires that
	binary.Read(f, binary.LittleEndian, &header)

	if err = header.Validate(); err != nil {
		return nil, err // header is invalid
	}

	// NOTE: explicit file close, see above
	f.Close()

	/* we now have a valid .way file, but we leave decoding it to the GetAccs() function */

	possibilities := getPossibleWayFormats()
	found := possibilities[header.FileType]
	// a key not found will return the value type's default value, in this case nil
	if found == nil {
		return nil, format.DetectionError(fmt.Sprintf("detector doesn't know what to do with format type #%d", header.FileType))
	}

	// prepopulate.
	err = found.PrepopulateFromFile(path, path_optional == nil)
	if err != nil {
		return nil, err
	}

	return found, nil
}

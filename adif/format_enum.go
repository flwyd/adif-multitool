// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package adif

import (
	"fmt"
	"strings"
)

const (
	// FormatADI is a Format of type ADI.
	FormatADI Format = iota
	// FormatADX is a Format of type ADX.
	FormatADX
	// FormatCSV is a Format of type CSV.
	FormatCSV
	// FormatJSON is a Format of type JSON.
	FormatJSON
)

const _FormatName = "ADIADXCSVJSON"

var _FormatNames = []string{
	_FormatName[0:3],
	_FormatName[3:6],
	_FormatName[6:9],
	_FormatName[9:13],
}

// FormatNames returns a list of possible string values of Format.
func FormatNames() []string {
	tmp := make([]string, len(_FormatNames))
	copy(tmp, _FormatNames)
	return tmp
}

var _FormatMap = map[Format]string{
	FormatADI:  _FormatName[0:3],
	FormatADX:  _FormatName[3:6],
	FormatCSV:  _FormatName[6:9],
	FormatJSON: _FormatName[9:13],
}

// String implements the Stringer interface.
func (x Format) String() string {
	if str, ok := _FormatMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Format(%d)", x)
}

var _FormatValue = map[string]Format{
	_FormatName[0:3]:                   FormatADI,
	strings.ToLower(_FormatName[0:3]):  FormatADI,
	_FormatName[3:6]:                   FormatADX,
	strings.ToLower(_FormatName[3:6]):  FormatADX,
	_FormatName[6:9]:                   FormatCSV,
	strings.ToLower(_FormatName[6:9]):  FormatCSV,
	_FormatName[9:13]:                  FormatJSON,
	strings.ToLower(_FormatName[9:13]): FormatJSON,
}

// ParseFormat attempts to convert a string to a Format.
func ParseFormat(name string) (Format, error) {
	if x, ok := _FormatValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _FormatValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Format(0), fmt.Errorf("%s is not a valid Format, try [%s]", name, strings.Join(_FormatNames, ", "))
}

// Set implements the Golang flag.Value interface func.
func (x *Format) Set(val string) error {
	v, err := ParseFormat(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *Format) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *Format) Type() string {
	return "Format"
}

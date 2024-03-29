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
	// SeparatorEmpty is a Separator of type Empty.
	SeparatorEmpty Separator = iota
	// SeparatorSpace is a Separator of type Space.
	SeparatorSpace
	// SeparatorTab is a Separator of type Tab.
	SeparatorTab
	// SeparatorNewline is a Separator of type Newline.
	SeparatorNewline
	// Separator2Newline is a Separator of type 2Newline.
	Separator2Newline
	// SeparatorCrlf is a Separator of type Crlf.
	SeparatorCrlf
	// Separator2Crlf is a Separator of type 2Crlf.
	Separator2Crlf
)

var ErrInvalidSeparator = fmt.Errorf("not a valid Separator, try [%s]", strings.Join(_SeparatorNames, ", "))

const _SeparatorName = "emptyspacetabnewline2newlinecrlf2crlf"

var _SeparatorNames = []string{
	_SeparatorName[0:5],
	_SeparatorName[5:10],
	_SeparatorName[10:13],
	_SeparatorName[13:20],
	_SeparatorName[20:28],
	_SeparatorName[28:32],
	_SeparatorName[32:37],
}

// SeparatorNames returns a list of possible string values of Separator.
func SeparatorNames() []string {
	tmp := make([]string, len(_SeparatorNames))
	copy(tmp, _SeparatorNames)
	return tmp
}

var _SeparatorMap = map[Separator]string{
	SeparatorEmpty:    _SeparatorName[0:5],
	SeparatorSpace:    _SeparatorName[5:10],
	SeparatorTab:      _SeparatorName[10:13],
	SeparatorNewline:  _SeparatorName[13:20],
	Separator2Newline: _SeparatorName[20:28],
	SeparatorCrlf:     _SeparatorName[28:32],
	Separator2Crlf:    _SeparatorName[32:37],
}

// String implements the Stringer interface.
func (x Separator) String() string {
	if str, ok := _SeparatorMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Separator(%d)", x)
}

var _SeparatorValue = map[string]Separator{
	_SeparatorName[0:5]:                    SeparatorEmpty,
	strings.ToLower(_SeparatorName[0:5]):   SeparatorEmpty,
	_SeparatorName[5:10]:                   SeparatorSpace,
	strings.ToLower(_SeparatorName[5:10]):  SeparatorSpace,
	_SeparatorName[10:13]:                  SeparatorTab,
	strings.ToLower(_SeparatorName[10:13]): SeparatorTab,
	_SeparatorName[13:20]:                  SeparatorNewline,
	strings.ToLower(_SeparatorName[13:20]): SeparatorNewline,
	_SeparatorName[20:28]:                  Separator2Newline,
	strings.ToLower(_SeparatorName[20:28]): Separator2Newline,
	_SeparatorName[28:32]:                  SeparatorCrlf,
	strings.ToLower(_SeparatorName[28:32]): SeparatorCrlf,
	_SeparatorName[32:37]:                  Separator2Crlf,
	strings.ToLower(_SeparatorName[32:37]): Separator2Crlf,
}

// ParseSeparator attempts to convert a string to a Separator.
func ParseSeparator(name string) (Separator, error) {
	if x, ok := _SeparatorValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _SeparatorValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Separator(0), fmt.Errorf("%s is %w", name, ErrInvalidSeparator)
}

// Set implements the Golang flag.Value interface func.
func (x *Separator) Set(val string) error {
	v, err := ParseSeparator(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *Separator) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *Separator) Type() string {
	return "Separator"
}

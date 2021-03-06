// Code generated by github.com/infraboard/mcube
// DO NOT EDIT

package geoip

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	enumDBFileContentTypeShowMap = map[DBFileContentType]string{
		IPv4Content:     "ipv4",
		LocationContent: "location",
	}

	enumDBFileContentTypeIDMap = map[string]DBFileContentType{
		"ipv4":     IPv4Content,
		"location": LocationContent,
	}
)

// ParseDBFileContentType Parse DBFileContentType from string
func ParseDBFileContentType(str string) (DBFileContentType, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := enumDBFileContentTypeIDMap[key]
	if !ok {
		return 0, fmt.Errorf("unknown Status: %s", str)
	}

	return v, nil
}

// Is todo
func (t DBFileContentType) Is(target DBFileContentType) bool {
	return t == target
}

// String stringer
func (t DBFileContentType) String() string {
	v, ok := enumDBFileContentTypeShowMap[t]
	if !ok {
		return "unknown"
	}

	return v
}

// MarshalJSON todo
func (t DBFileContentType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(t.String())
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *DBFileContentType) UnmarshalJSON(b []byte) error {
	ins, err := ParseDBFileContentType(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

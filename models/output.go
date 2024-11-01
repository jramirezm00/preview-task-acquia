package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FinalResponseItem struct {
	Title   string
	Planets []string
}

type FinalResponse []FinalResponseItem

// custom marshaller in order to respect key order of the titles
func (fr FinalResponse) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, mi := range fr {
		b, err := json.Marshal(&mi.Planets)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprintf("%v", mi.Title)))
		buf.Write(b)
		if i < len(fr)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

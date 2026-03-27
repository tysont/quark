// ABOUTME: Serialization utilities for converting structs to and
// ABOUTME: from bytes using Go's gob encoding.
package quark

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func encode(o interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(o)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	return b, nil
}

func decode(t reflect.Type, b []byte) (interface{}, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	o := reflect.New(t.Elem()).Interface()
	err := decoder.Decode(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
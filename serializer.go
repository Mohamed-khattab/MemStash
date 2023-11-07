package main

import (
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

//// newWriter creates a new Writer instance with the provided io.Writer.
func newWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "error":
		return v.marshalError()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

// MarshalArray serializes an array of values into bytes.
func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')
	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes

}

// MarshalBulk marshals the Value into a bulk format.
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// Write writes the given value to the writer.
func (w *Writer) Write(v Value) error {
	_, err := w.writer.Write(v.Marshal())
	return err
}

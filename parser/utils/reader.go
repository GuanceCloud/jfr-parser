package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type JFRReader struct {
	r    io.Reader
	buf  []byte
	pos  int
	size int
	err  error
}

func NewDataReader(r io.Reader) *JFRReader {
	return &JFRReader{
		r: r,
	}
}

func (d *JFRReader) FillTo(size int) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}
	if len(d.buf) < size {
		d.buf = append(d.buf, make([]byte, size-len(d.buf))...)
	}
	n, err = io.ReadFull(d.r, d.buf[d.size:])
	d.size += n
	return n, err
}

func (d *JFRReader) Buf() []byte {
	return d.buf
}

func (d *JFRReader) Pos() int {
	return d.pos
}

func (d *JFRReader) Unread() int {
	return d.size - d.pos
}

func (d *JFRReader) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	var errFill error
	if d.size-d.pos < len(p) {
		_, errFill = d.FillTo(d.pos + len(p))
	}

	n = copy(p, d.buf[d.pos:d.size])
	d.pos += n
	if n > 0 {
		return n, nil
	}
	d.err = errFill
	return 0, d.err
}

func (d *JFRReader) BigEndianRead(x any) error {
	if err := binary.Read(d, binary.BigEndian, x); err != nil {
		rv := reflect.ValueOf(x)
		for rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		return fmt.Errorf("unable to read next %T: %w", rv.Interface(), err)
	}
	return nil
}

func (d *JFRReader) ReadInt8() (b int8, err error) {
	err = d.BigEndianRead(&b)
	return
}

func (d *JFRReader) ReadByte() (b byte, err error) {
	err = d.BigEndianRead(&b)
	return
}

func (d *JFRReader) ReadBool() (bool, error) {
	b, err := d.ReadByte()
	if err != nil {
		return false, fmt.Errorf("unable to read next boolean: %w", err)
	}
	return b != 0, nil
}

func (d *JFRReader) ReadShort() (n int16, err error) {
	err = d.BigEndianRead(&n)
	return
}

func (d *JFRReader) ReadUint16() (uint16, error) {
	var u16 uint16
	if err := binary.Read(d, binary.BigEndian, &u16); err != nil {
		return 0, fmt.Errorf("unable to read next uint16: %w", err)
	}
	return u16, nil
}

func (d *JFRReader) ReadUint32() (u32 uint32, err error) {
	err = d.BigEndianRead(&u32)
	return
}

func (d *JFRReader) ReadInt32() (i32 int32, err error) {
	err = d.BigEndianRead(&i32)
	return
}

func (d *JFRReader) ReadLong() (int64, error) {
	var i64 int64
	if err := binary.Read(d, binary.BigEndian, &i64); err != nil {
		return 0, fmt.Errorf("unable to read next int64: %w", err)
	}
	return i64, nil
}

func (d *JFRReader) ReadFloat32() (f32 float32, err error) {
	err = d.BigEndianRead(&f32)
	return
}

func (d *JFRReader) ReadDouble() (f64 float64, err error) {
	err = d.BigEndianRead(&f64)
	return
}

func (d *JFRReader) ReadString() (string, error) {
	length, err := d.ReadUint16()
	if err != nil {
		return "", fmt.Errorf("unable to read string length: %w", err)
	}

	byteSlice := make([]byte, length)
	if _, err = io.ReadFull(d, byteSlice); err != nil {
		return "", fmt.Errorf("unable to read string data: %w", err)
	}
	return string(byteSlice), nil
}

func (d *JFRReader) Skip(n int) (int, error) {
	if d.err != nil {
		return 0, nil
	}
	var errFill error
	if d.Unread() < n {
		_, errFill = d.FillTo(d.pos + n)
	}

	if d.Unread() < n {
		skipped := d.Unread()
		d.pos += skipped
		d.err = errFill
		return skipped, d.err
	}

	d.pos += n
	return n, nil
}

// ReadAt reads len(p) bytes into p starting at offset off in the
// underlying input source, ReadAt doesn't forward current reading offset.
// When ReadAt returns n < len(p), it returns a non-nil error
// explaining why more bytes were not returned.
func (d *JFRReader) ReadAt(p []byte, offset int64) (n int, err error) {
	var errFill error
	if d.size-int(offset) < len(p) {
		_, errFill = d.FillTo(int(offset) + len(p))
	}

	if offset < int64(d.size) {
		n = copy(p, d.buf[offset:d.size])
	}

	if n < len(p) {
		if errFill != nil {
			return n, errFill
		}
		return n, io.ErrUnexpectedEOF
	}
	return n, nil
}

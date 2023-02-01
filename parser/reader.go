package parser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	reader2 "github.com/GuanceCloud/jfr-parser/reader"
)

type Reader interface {
	Boolean() (bool, error)
	Byte() (int8, error)
	Short() (int16, error)
	Char() (uint16, error)
	Int() (int32, error)
	Long() (int64, error)
	Float() (float32, error)
	Double() (float64, error)
	String(pool *CPool) (string, error)

	reader2.VarReader

	// TODO: Support arrays
}

type InputReader interface {
	io.Reader
	io.ByteReader
}

type reader struct {
	InputReader
	varR reader2.VarReader
}

func NewReader(r InputReader, compressed bool) Reader {
	var varR reader2.VarReader
	if compressed {
		varR = reader2.NewCompressed(r)
	} else {
		varR = reader2.NewUncompressed(r)
	}
	return reader{
		InputReader: r,
		varR:        varR,
	}
}

func (r reader) Boolean() (bool, error) {
	var n int8
	err := binary.Read(r, binary.BigEndian, &n)
	if n == 0 {
		return false, err
	}
	return true, err
}

func (r reader) Byte() (int8, error) {
	var n int8
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Short() (int16, error) {
	return reader2.Short(r)
}

func (r reader) Char() (uint16, error) {
	var n uint16
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Int() (int32, error) {
	return reader2.Int(r)
}

func (r reader) Long() (int64, error) {
	return reader2.Long(r)
}

func (r reader) Float() (float32, error) {
	var n float32
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Double() (float64, error) {
	var n float64
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

// TODO: Should we differentiate between null and empty?
func (r reader) String(pool *CPool) (string, error) {
	enc, err := r.Byte()
	if err != nil {
		return "", err
	}
	switch enc {
	case 0:
		return "", nil
	case 1:
		return "", nil
	case 2: // constant pool reference
		idx, err := r.VarLong()
		if err != nil {
			fmt.Printf("get constant refrence idx fail: %s\n", err)
			return "", err
		}
		if pool == nil {
			return "", errors.New("the string constant pool is nil")
		}
		v, ok := pool.Pool[int(idx)]
		if !ok {
			return "", fmt.Errorf("string not found in the pool")
		}
		str, ok := v.(*String)
		if !ok {
			return "", fmt.Errorf("not type of parser.String")
		}
		return string(*str), nil
	case 3, 4, 5:
		return r.utf8()
	default:
		// TODO
		return "", fmt.Errorf("Unsupported string type :%d", enc)
	}
}

func (r reader) VarShort() (int16, error) {
	return r.varR.VarShort()
}

func (r reader) VarInt() (int32, error) {
	return r.varR.VarInt()
}

func (r reader) VarLong() (int64, error) {
	return r.varR.VarLong()
}

func (r reader) utf8() (string, error) {
	n, err := r.varR.VarInt()
	if err != nil {
		return "", nil
	}
	// TODO: make sure n is reasonable
	b := make([]byte, n)
	_, err = io.ReadFull(r, b)
	return string(b), err
}

package parser

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChunkReader(t *testing.T) {
	buf := make([]byte, 1<<20)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		t.Fatal(err)
	}

	cr := NewChunkReader(bytes.NewBuffer(buf))

	all, err := io.ReadAll(cr)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("pos: %d, size: %d\n", cr.pos, cr.size)
	assert.Equal(t, 0, bytes.Compare(buf, all))
	assert.Equal(t, 0, bytes.Compare(buf, cr.buf[:cr.size]))
	assert.Equal(t, 0, cr.Unread())

	out := make([]byte, 345)
	n, err := cr.ReadAt(out, 12345)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(out), n)
	assert.Equal(t, 0, bytes.Compare(out, buf[12345:12345+len(out)]))

	n, err = cr.ReadAt(out[:5], int64(cr.size)-5)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, n)
	assert.Equal(t, 0, bytes.Compare(out[:5], buf[len(buf)-5:]))
}

func TestChunkReader_Skip(t *testing.T) {
	buf := make([]byte, 1<<20)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		t.Fatal(err)
	}

	cr := NewChunkReader(bytes.NewBuffer(buf))

	_, err := cr.Skip(11111)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cr.Skip(11111 * 2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cr.Skip(11111 * 3)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cr.Skip(11111 * 5)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cr.Skip(11111 * 8)
	if err != nil {
		t.Fatal(err)
	}

	all, err := io.ReadAll(cr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(buf)-11111*(1+2+3+5+8), len(all))
	assert.Equal(t, 0, bytes.Compare(all, buf[11111*(1+2+3+5+8):]))

	t.Logf("cr.err: %v", cr.err)
	n, err := cr.FillTo(cr.size + 5)
	assert.NotNil(t, err)
	assert.Equal(t, 0, n)
}

func TestChunkReader_FillTo(t *testing.T) {
	buf := make([]byte, 1000)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		t.Fatal(err)
	}

	cr := NewChunkReader(bytes.NewBuffer(buf))

	n, err := cr.FillTo(8)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 8, n)
	assert.Equal(t, 8, cr.size)

	out := make([]byte, 8)
	_, err = io.ReadFull(cr, out)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, bytes.Compare(out, buf[:8]))

	n, err = cr.Skip(456)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 456, n)

	n, err = cr.FillTo(916)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fillto res: %d", n)

	n, err = cr.Skip(800)
	t.Log(err)
	assert.NotNil(t, err)
	assert.Equal(t, 1000-8-456, n)
}

func TestChunkReader_ReadAt(t *testing.T) {
	buf := make([]byte, 1<<20)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		t.Fatal(err)
	}

	cr := NewChunkReader(bytes.NewBuffer(buf))

	out := make([]byte, 996)

	n, err := cr.ReadAt(out, 20240)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(out), n)
	assert.Equal(t, 0, bytes.Compare(out, buf[20240:20240+len(out)]))
	assert.Equal(t, 0, cr.pos)

	n, err = cr.ReadAt(out, int64(len(buf))-365)
	t.Log(err)
	assert.NotNil(t, err)
	assert.Equal(t, 365, n)
	assert.Equal(t, 0, bytes.Compare(out[:n], buf[len(buf)-365:]))
}

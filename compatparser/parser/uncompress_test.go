package parser

import (
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func TestUncompressZSTD(t *testing.T) {
	const plainText = "hello world\n\n"

	var compressed bytes.Buffer
	zw, err := zstd.NewWriter(&compressed)
	if err != nil {
		t.Fatalf("unable to create zstd writer: %s", err)
	}
	if _, err := zw.Write([]byte(plainText)); err != nil {
		t.Fatalf("unable to write zstd payload: %s", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("unable to close zstd writer: %s", err)
	}

	r, err := Uncompress(bytes.NewReader(compressed.Bytes()))
	if err != nil {
		t.Fatalf("unable to uncompress zstd payload: %s", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read uncompressed data err: %s", err)
	}
	if string(data) != plainText {
		t.Fatalf("uncompressed data incorrect, expect [%s], [%s] found", plainText, string(data))
	}
}

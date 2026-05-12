package parser

import (
	"bytes"
	"testing"
)

func TestReaderStringCharArray(t *testing.T) {
	rd := NewReader(bytes.NewReader([]byte{4, 2, 0xb5, 0x01, 's'}), true)

	got, err := rd.String(nil)
	if err != nil {
		t.Fatalf("String() error = %v", err)
	}
	if got != "µs" {
		t.Fatalf("String() = %q, want %q", got, "µs")
	}
}

func TestReaderStringLatin1(t *testing.T) {
	rd := NewReader(bytes.NewReader([]byte{5, 1, 0xe9}), true)

	got, err := rd.String(nil)
	if err != nil {
		t.Fatalf("String() error = %v", err)
	}
	if got != "é" {
		t.Fatalf("String() = %q, want %q", got, "é")
	}
}

package parser

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type testCase struct {
	name string
	file []byte
}

var (
	plainText = "hello world\n\n"

	testCases = []testCase{
		{
			name: "zip",
			file: []byte{80, 75, 3, 4, 10, 0, 0, 0, 0, 0, 123, 11, 47, 86, 221, 82, 167, 119, 13, 0, 0, 0, 13, 0, 0, 0, 5, 0, 28, 0, 50, 46, 116, 120, 116, 85, 84, 9, 0, 3, 25, 230, 194, 99, 26, 230, 194, 99, 117, 120, 11, 0, 1, 4, 245, 1, 0, 0, 4, 20, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 10, 10, 80, 75, 1, 2, 30, 3, 10, 0, 0, 0, 0, 0, 123, 11, 47, 86, 221, 82, 167, 119, 13, 0, 0, 0, 13, 0, 0, 0, 5, 0, 24, 0, 0, 0, 0, 0, 1, 0, 0, 0, 164, 129, 0, 0, 0, 0, 50, 46, 116, 120, 116, 85, 84, 5, 0, 3, 25, 230, 194, 99, 117, 120, 11, 0, 1, 4, 245, 1, 0, 0, 4, 20, 0, 0, 0, 80, 75, 5, 6, 0, 0, 0, 0, 1, 0, 1, 0, 75, 0, 0, 0, 76, 0, 0, 0, 0, 0},
		},
		{
			name: "gzip",
			file: []byte{31, 139, 8, 8, 25, 230, 194, 99, 0, 3, 50, 46, 116, 120, 116, 0, 203, 72, 205, 201, 201, 87, 40, 207, 47, 202, 73, 225, 226, 2, 0, 221, 82, 167, 119, 13, 0, 0, 0},
		},
		{
			name: "lz4",
			file: []byte{4, 34, 77, 24, 100, 64, 167, 13, 0, 0, 128, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 10, 10, 0, 0, 0, 0, 157, 89, 174, 11},
		},
	}
)

func TestUncompress(t *testing.T) {

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := Decompress(bytes.NewReader(tc.file))
			if err != nil {
				t.Fatalf("unable to uncompress file: %s", err)
			}

			defer r.Close()

			data, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("read uncompress data err: %s", err)
			}

			fmt.Println(string(data))

			if string(data) != plainText {
				t.Fatalf("uncompress zip file incorrect, expect [%s], [%s] found", plainText, string(data))
			}
		})
	}
}

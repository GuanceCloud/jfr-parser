package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func testParseFile(name string) ([]Chunk, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open jfr file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

func TestParseUncompressed(t *testing.T) {
	chunks, err := testParseFile("./testdata/ddtrace.jfr")
	if err != nil {
		t.Fatalf("Unable to parse jfr file: %s", err)
	}
	fmt.Println("chunks length: ", len(chunks))
}

func TestParseZip(t *testing.T) {
	chunks, err := testParseFile("./testdata/ddtrace.jfr.zip")
	if err != nil {
		t.Fatalf("Unable to parse jfr file: %s", err)
	}
	fmt.Println("chunks length: ", len(chunks))
}

func TestParseLZ4(t *testing.T) {

	chunks, err := testParseFile("./testdata/ddtrace.jfr.lz4")
	if err != nil {
		t.Fatalf("Unable to parse jfr file: %s", err)
	}

	fmt.Println("chunks length:", len(chunks))
}

func TestParse(t *testing.T) {
	jfr, err := os.Open("./testdata/example.jfr.gz")
	if err != nil {
		t.Fatalf("Unable to read JFR file: %s", err)
	}
	expectedJson, err := readGzipFile("./testdata/example_parsed.json.gz")
	if err != nil {
		t.Fatalf("Unable to read example_parsd.json")
	}
	chunks, err := Parse(jfr)
	if err != nil {
		t.Fatalf("Failed to parse JFR: %s", err)
		return
	}
	actualJson, _ := json.Marshal(chunks)
	if !bytes.Equal(expectedJson, actualJson) {
		t.Fatalf("Failed to parse JFR: %s", err)
		return
	}
}

func BenchmarkParse(b *testing.B) {
	jfr, err := os.Open("./testdata/example.jfr.gz")
	if err != nil {
		b.Fatalf("Unable to read JFR file: %s", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := Parse(jfr)
		if err != nil {
			b.Fatalf("Unable to parse JFR file: %s", err)
		}
	}
}

func readGzipFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

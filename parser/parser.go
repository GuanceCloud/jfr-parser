package parser

import (
	"fmt"
	"io"
)

func Parse(r io.Reader) ([]Chunk, error) {
	rc, err := Uncompress(r)
	if err != nil {
		return nil, fmt.Errorf("unable to decompress input stream: %w", err)
	}
	defer rc.Close()
	return ParseWithOptions(rc, &ChunkParseOptions{})
}

func ParseWithOptions(r io.Reader, options *ChunkParseOptions) ([]Chunk, error) {
	var chunks []Chunk
	for {
		var chunk Chunk
		err := chunk.Parse(r, options)
		if err == io.EOF {
			return chunks, nil
		}
		if err != nil {
			return chunks, fmt.Errorf("unable to parse chunk: %w", err)
		}
		chunks = append(chunks, chunk)
	}
}

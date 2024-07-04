package parser

import (
	"bytes"
	"fmt"
	"io"
)

const (
	MetadataEventType     = 0
	ConstantPoolEventType = 1

	EventSuperType = "jdk.jfr.Event"
)

var magic = []byte{'F', 'L', 'R', 0}

type Version struct {
	Major uint16
	Minor uint16
}

type CPool struct {
	Pool     map[int]ParseResolvable
	resolved bool
}
type ClassMap map[int]*ClassMetadata
type PoolMap map[int]*CPool

type Chunk struct {
	Header      Header
	Metadata    MetadataEvent
	Checkpoints []CheckpointEvent
	Events      []Event
}

type ChunkParseOptions struct {
	CPoolProcessor func(meta *ClassMetadata, cpool *CPool)
}

func (c *Chunk) Parse(r io.Reader, options *ChunkParseOptions) (err error) {
	buf := make([]byte, len(magic))
	if _, err = io.ReadFull(r, buf); err != nil {
		if err == io.EOF {
			return err
		}
		return fmt.Errorf("unable to read chunk's header: %w", err)
	}

	// TODO magic header
	for i, r := range magic {
		if r != buf[i] {
			return fmt.Errorf("unexpected magic header %v expected, %v found", magic, buf)
		}
	}
	if _, err = io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("unable to read format version: %w", err)
	}

	// TODO Check supported major / minor

	buf = make([]byte, headerSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("unable to read chunk header: %w", err)
	}
	if err := c.Header.Parse(NewReader(bytes.NewReader(buf), false)); err != nil {
		return fmt.Errorf("unable to parse chunk header: %w", err)
	}
	c.Header.ChunkSize -= headerSize + 8
	c.Header.MetadataOffset -= headerSize + 8
	c.Header.ConstantPoolOffset -= headerSize + 8
	useCompression := c.Header.Features&1 == 1
	// TODO: assert c.Header.ChunkSize is small enough
	buf = make([]byte, c.Header.ChunkSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("unable to read chunk contents: %w", err)
	}

	br := bytes.NewReader(buf)
	rd := NewReader(br, useCompression)
	pointer := int64(0)
	events := make(map[int64]int32)

	// Parse metadata
	if _, err := br.Seek(c.Header.MetadataOffset, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	metadataSize, err := rd.VarInt()
	if err != nil {
		return fmt.Errorf("unable to parse chunk metadata size: %w", err)
	}
	events[c.Header.MetadataOffset] = metadataSize
	if err := c.Metadata.Parse(rd); err != nil {
		return fmt.Errorf("unable to parse chunk metadata: %w", err)
	}
	classes := buildClasses(c.Metadata)

	// Parse checkpoint event(s)
	if _, err := br.Seek(c.Header.ConstantPoolOffset, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	checkpointsSize := int32(0)
	cpools := make(PoolMap)
	delta := int64(0)
	for {
		size, err := rd.VarInt()
		if err != nil {
			return fmt.Errorf("unable to parse checkpoint event size: %w", err)
		}
		events[c.Header.ConstantPoolOffset+delta] = size
		checkpointsSize += size
		var cp CheckpointEvent
		if err := cp.Parse(rd, classes, cpools); err != nil {
			return fmt.Errorf("unable to parse checkpoint event: %w", err)
		}
		c.Checkpoints = append(c.Checkpoints, cp)
		if cp.Delta == 0 {
			break
		}
		delta += cp.Delta
		if _, err := br.Seek(c.Header.ConstantPoolOffset+delta, io.SeekStart); err != nil {
			return fmt.Errorf("unable to seek reader: %w", err)
		}
	}

	if options.CPoolProcessor != nil {
		for classID, pool := range cpools {
			options.CPoolProcessor(classes[classID], pool)
		}
	}

	// Second pass over constant pools: resolve constants
	for classID := range cpools {
		if err := ResolveConstants(classes, cpools, classID); err != nil {
			return err
		}
	}

	// Parse the rest of events
	if _, err := br.Seek(pointer, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	for pointer != c.Header.ChunkSize {
		if size, ok := events[pointer]; ok {
			pointer += int64(size)
		} else {
			if _, err := br.Seek(pointer, io.SeekStart); err != nil {
				return fmt.Errorf("unable to seek to position %d: %w", pointer, err)
			}
			size, err := rd.VarInt()
			if err != nil {
				return fmt.Errorf("unable to parse event size: %w", err)
			}
			if size == 0 {
				return fmt.Errorf("found event with invalid size (0)")
			}
			events[pointer] = size
			e, err := ParseEvent(rd, classes, cpools)
			if err != nil {
				return fmt.Errorf("unable to parse event: %w", err)
			}
			if e != nil {
				c.Events = append(c.Events, e)
			}
			pointer += int64(size)
		}
	}
	return nil
}

func buildClasses(metadata MetadataEvent) ClassMap {
	classes := make(ClassMap)
	for _, class := range metadata.Root.Metadata.Classes {
		classes[int(class.ID)] = class
	}
	return classes
}

func ResolveConstants(classes ClassMap, cpools PoolMap, classID int) (err error) {
	cpool, ok := cpools[classID]
	if !ok {
		// Non-existent constant pool references seem to be used to mark no value
		return nil
	}
	if cpool.resolved {
		return nil
	}
	cpool.resolved = true
	for _, t := range cpool.Pool {
		if err := t.Resolve(classes, cpools); err != nil {
			return fmt.Errorf("unable to resolve constants: %w", err)
		}
	}
	return nil
}

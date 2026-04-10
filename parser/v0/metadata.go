package v0

import (
	"encoding/binary"
	"fmt"
	"github.com/grafana/jfr-parser/common/units"
)

const (
	ByteSize    = 1
	BooleanSize = 1
	ShortSize   = 2
	CharSize    = 2
	IntegerSize = 4
	LongSize    = 8
	FloatSize   = 4
	DoubleSize  = 8
)

const (
	MetadataEventTypeIndex   = 0
	CheckPointEventTypeIndex = 1
	LostEventTypeIndex       = 2
)

type Offset struct {
	offset      int64
	offsetLimit int64
}

func NewOffset(data []byte, startOffset int64) *Offset {
	size := binary.BigEndian.Uint32(data[startOffset : startOffset+IntegerSize])

	offsetLimit := startOffset + int64(size)

	return &Offset{
		offset:      startOffset + IntegerSize,
		offsetLimit: offsetLimit,
	}
}

func (o *Offset) Incr(n int64) error {
	if o.offset+n > o.offsetLimit {
		return fmt.Errorf("offset out of max offset limit: %d", o.offset+n)
	}
	o.offset += n
	return nil
}

func (o *Offset) GetAndIncr(n int64) (int64, error) {
	cur := o.offset
	if err := o.Incr(n); err != nil {
		return cur, fmt.Errorf("unable to incr offset: %w", err)
	}
	return cur, nil
}

type ChunkStruct struct {
	MetadataOffset  int64
	BodyStartOffset int64
	ChunkSize       int64
}

type ChunkMetadata struct {
	*ChunkStruct
	Producers          []*ProducerDescriptor
	StartTimeUnixNano  int64
	EndTimeUnixNano    int64
	StartTicks         int64
	TicksPerNano       float64
	PreviousCheckPoint int64
	UnitTick           *units.Unit
	ConstantMapPool    ConstantMapPool
}

func (c *ChunkMetadata) AsUnixNanoTimeStamp(endTicks int64) int64 {
	return c.StartTimeUnixNano + int64(float64(endTicks-c.StartTicks)/c.TicksPerNano)
}

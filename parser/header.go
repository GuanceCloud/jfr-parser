package parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

const headerSize = 60

type Header struct {
	ChunkSize          int64
	ConstantPoolOffset int64
	MetadataOffset     int64
	StartTimeNanos     int64
	DurationNanos      int64
	StartTicks         int64
	TicksPerSecond     int64
	Features           int32
}

func (h *Header) Parse(rd Reader) (err error) {
	h.ChunkSize, _ = rd.Long()
	h.ConstantPoolOffset, _ = rd.Long()
	h.MetadataOffset, _ = rd.Long()
	h.StartTimeNanos, _ = rd.Long()
	h.DurationNanos, _ = rd.Long()
	h.StartTicks, _ = rd.Long()
	h.TicksPerSecond, _ = rd.Long()
	h.Features, err = rd.Int()
	return err
}

func (p *Parser) readChunkHeader(pos int) error {
	if p.pos+chunkHeaderSize > len(p.buf) {
		return io.ErrUnexpectedEOF
	}

	p.pos = pos
	h := ChunkHeader{}
	h.Features = binary.BigEndian.Uint32(p.buf[pos+64:])
	h.Magic = binary.BigEndian.Uint32(p.buf[pos:])
	h.Version = binary.BigEndian.Uint32(p.buf[pos+4:])
	h.Size = int(binary.BigEndian.Uint64(p.buf[pos+8:]))
	h.OffsetConstantPool = int(binary.BigEndian.Uint64(p.buf[pos+16:]))
	h.OffsetMeta = int(binary.BigEndian.Uint64(p.buf[pos+24:]))
	h.StartNanos = binary.BigEndian.Uint64(p.buf[pos+32:])
	h.DurationNanos = binary.BigEndian.Uint64(p.buf[pos+40:])
	h.StartTicks = binary.BigEndian.Uint64(p.buf[pos+48:])
	h.TicksPerSecond = binary.BigEndian.Uint64(p.buf[pos+56:])
	if h.Magic != chunkMagic {
		return fmt.Errorf("invalid chunk magic: %x", h.Magic)
	}
	if h.Version < 0x20000 || h.Version > 0x2ffff {
		return fmt.Errorf("unknown version %x", h.Version)
	}
	if h.OffsetConstantPool <= 0 || h.OffsetMeta <= 0 {
		return fmt.Errorf("invalid offsets: cp %d meta %d", h.OffsetConstantPool, h.OffsetMeta)
	}
	if h.Size <= 0 {
		return fmt.Errorf("invalid size: %d", h.Size)
	}
	if p.options.ChunkSizeLimit > 0 && h.Size > p.options.ChunkSizeLimit {
		return fmt.Errorf("chunk size %d exceeds limit %d", h.Size, p.options.ChunkSizeLimit)
	}
	p.header = h
	p.chunkEnd = pos + h.Size
	return nil
}

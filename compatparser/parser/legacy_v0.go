package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf16"
)

const (
	legacyV0MetadataEventTypeIndex   = 0
	legacyV0CheckPointEventTypeIndex = 1

	legacyV0ByteSize    = 1
	legacyV0IntegerSize = 4
	legacyV0LongSize    = 8

	legacyV0ContentTypePoolNone    = 0
	legacyV0ContentTypeMemory      = 1
	legacyV0ContentTypeEpochMillis = 2
	legacyV0ContentTypeMillis      = 3
	legacyV0ContentTypeNanos       = 4
	legacyV0ContentTypeTicks       = 5
	legacyV0ContentTypeAddress     = 6
	legacyV0ContentTypeThreadID    = 7
	legacyV0ContentTypeJavaThread  = 8
	legacyV0ContentTypeStackTrace  = 9
)

type legacyV0Reader struct {
	r    io.Reader
	buf  []byte
	pos  int
	size int
	err  error
}

func newLegacyV0Reader(r io.Reader) *legacyV0Reader {
	return &legacyV0Reader{r: r}
}

func (d *legacyV0Reader) FillTo(size int) (int, error) {
	if d.err != nil {
		return 0, d.err
	}
	if len(d.buf) < size {
		d.buf = append(d.buf, make([]byte, size-len(d.buf))...)
	}
	n, err := io.ReadFull(d.r, d.buf[d.size:size])
	d.size += n
	return n, err
}

func (d *legacyV0Reader) Buf() []byte { return d.buf[:d.size] }
func (d *legacyV0Reader) Pos() int    { return d.pos }
func (d *legacyV0Reader) Unread() int { return d.size - d.pos }

func (d *legacyV0Reader) Read(p []byte) (int, error) {
	if d.err != nil {
		return 0, d.err
	}
	var errFill error
	if d.size-d.pos < len(p) {
		_, errFill = d.FillTo(d.pos + len(p))
	}
	n := copy(p, d.buf[d.pos:d.size])
	d.pos += n
	if n > 0 {
		return n, nil
	}
	d.err = errFill
	return 0, d.err
}

func (d *legacyV0Reader) ReadAt(p []byte, offset int64) (int, error) {
	var errFill error
	if d.size-int(offset) < len(p) {
		_, errFill = d.FillTo(int(offset) + len(p))
	}
	n := 0
	if offset < int64(d.size) {
		n = copy(p, d.buf[offset:d.size])
	}
	if n < len(p) {
		if errFill != nil {
			return n, errFill
		}
		return n, io.ErrUnexpectedEOF
	}
	return n, nil
}

func (d *legacyV0Reader) Skip(n int) (int, error) {
	if d.err != nil {
		return 0, d.err
	}
	var errFill error
	if d.Unread() < n {
		_, errFill = d.FillTo(d.pos + n)
	}
	if d.Unread() < n {
		skipped := d.Unread()
		d.pos += skipped
		d.err = errFill
		return skipped, d.err
	}
	d.pos += n
	return n, nil
}

func (d *legacyV0Reader) readBigEndian(v any) error {
	if err := binary.Read(d, binary.BigEndian, v); err != nil {
		return err
	}
	return nil
}

func (d *legacyV0Reader) ReadByte() (byte, error) {
	var b byte
	return b, d.readBigEndian(&b)
}

func (d *legacyV0Reader) ReadBool() (bool, error) {
	b, err := d.ReadByte()
	return b != 0, err
}

func (d *legacyV0Reader) ReadInt8() (int8, error) {
	var n int8
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadUint16() (uint16, error) {
	var n uint16
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadShort() (int16, error) {
	var n int16
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadInt32() (int32, error) {
	var n int32
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadUint32() (uint32, error) {
	var n uint32
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadLong() (int64, error) {
	var n int64
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadFloat32() (float32, error) {
	var n float32
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadDouble() (float64, error) {
	var n float64
	return n, d.readBigEndian(&n)
}

func (d *legacyV0Reader) ReadString() (string, error) {
	length, err := d.ReadUint16()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(d, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (d *legacyV0Reader) ReadUTF16String() (string, error) {
	charCount, err := d.ReadInt32()
	if err != nil {
		return "", err
	}
	if charCount < 0 {
		return "", fmt.Errorf("illegal utf16 char count %d", charCount)
	}
	data := make([]byte, int(charCount)*2)
	if _, err := io.ReadFull(d, data); err != nil {
		return "", err
	}
	u16 := make([]uint16, 0, len(data)/2)
	for i := 0; i+1 < len(data); i += 2 {
		u16 = append(u16, binary.BigEndian.Uint16(data[i:i+2]))
	}
	return string(utf16.Decode(u16)), nil
}

type legacyV0DataType struct {
	Name    string
	Size    int
	Numeric bool
	Textual bool
}

func (d *legacyV0DataType) IsPrimitive() bool {
	switch d {
	case legacyV0Boolean, legacyV0Byte, legacyV0U1, legacyV0Short, legacyV0U2,
		legacyV0Integer, legacyV0U4, legacyV0Long, legacyV0U8, legacyV0Float, legacyV0Double,
		legacyV0UTF8, legacyV0String:
		return true
	default:
		return false
	}
}

var (
	legacyV0Boolean     = &legacyV0DataType{Name: "Boolean", Size: 1}
	legacyV0Byte        = &legacyV0DataType{Name: "Byte", Size: 1, Numeric: true}
	legacyV0U1          = &legacyV0DataType{Name: "U1", Size: 1, Numeric: true}
	legacyV0Short       = &legacyV0DataType{Name: "Short", Size: 2, Numeric: true}
	legacyV0U2          = &legacyV0DataType{Name: "U2", Size: 2, Numeric: true}
	legacyV0Integer     = &legacyV0DataType{Name: "Integer", Size: 4, Numeric: true}
	legacyV0U4          = &legacyV0DataType{Name: "U4", Size: 4, Numeric: true}
	legacyV0Long        = &legacyV0DataType{Name: "Long", Size: 8, Numeric: true}
	legacyV0U8          = &legacyV0DataType{Name: "U8", Size: 8, Numeric: true}
	legacyV0Float       = &legacyV0DataType{Name: "Float", Size: 4, Numeric: true}
	legacyV0Double      = &legacyV0DataType{Name: "Double", Size: 8, Numeric: true}
	legacyV0UTF8        = &legacyV0DataType{Name: "UTF8", Textual: true}
	legacyV0String      = &legacyV0DataType{Name: "String", Textual: true}
	legacyV0Array       = &legacyV0DataType{Name: "Array"}
	legacyV0Struct      = &legacyV0DataType{Name: "Struct"}
	legacyV0StructArray = &legacyV0DataType{Name: "StructArray"}
	legacyV0DataTypes   = []*legacyV0DataType{
		legacyV0Boolean, legacyV0Byte, legacyV0U1, legacyV0Short, legacyV0U2,
		legacyV0Integer, legacyV0U4, legacyV0Long, legacyV0U8, legacyV0Float,
		legacyV0Double, legacyV0UTF8, legacyV0String, legacyV0Array,
		legacyV0Struct, legacyV0StructArray,
	}
)

type legacyV0DataStructure struct {
	ValueDescriptors []*legacyV0ValueDescriptor
}

type legacyV0DataStructureList []*legacyV0DataStructure

func (d *legacyV0DataStructureList) Add(ds *legacyV0DataStructure) {
	*d = append(*d, ds)
}

func (d legacyV0DataStructureList) Get(i int64) *legacyV0DataStructure {
	return d[i]
}

type legacyV0ValueDescriptor struct {
	Identifier     string
	Name           string
	Description    string
	Transition     int
	DataType       *legacyV0DataType
	RelationalKey  string
	ContentType    int
	DataStructures *legacyV0DataStructureList
	StructureIndex int
}

func (v *legacyV0ValueDescriptor) ArrayElemDatatype() (*legacyV0DataType, error) {
	if v.DataType != legacyV0Array {
		return nil, fmt.Errorf("value descriptor [%s] is not an array", v.Identifier)
	}
	return legacyV0DataTypes[v.StructureIndex], nil
}

func (v *legacyV0ValueDescriptor) Children() ([]*legacyV0ValueDescriptor, error) {
	if v.DataType != legacyV0Struct && v.DataType != legacyV0StructArray {
		return nil, fmt.Errorf("value descriptor [%s] has no children", v.Identifier)
	}
	return v.DataStructures.Get(int64(v.StructureIndex)).ValueDescriptors, nil
}

func (v *legacyV0ValueDescriptor) InConstantPool() bool {
	if !v.DataType.Numeric {
		return false
	}
	switch v.ContentType {
	case legacyV0ContentTypeEpochMillis, legacyV0ContentTypeMillis, legacyV0ContentTypeNanos,
		legacyV0ContentTypeTicks, legacyV0ContentTypeAddress, legacyV0ContentTypeMemory,
		legacyV0ContentTypePoolNone:
		return false
	default:
		return true
	}
}

type legacyV0ProducerDescriptor struct {
	Name                   string
	Description            string
	URI                    string
	Identifier             int
	EventTypeDescriptors   []*legacyV0EventTypeDescriptor
	ContentTypeDescriptors []*legacyV0ContentTypeDescriptor
}

type legacyV0EventTypeDescriptor struct {
	Identifier        int
	Label             string
	HasStartTime      bool
	HasThread         bool
	CanHaveStacktrace bool
	IsRequestable     bool
	ValueDescriptors  []*legacyV0ValueDescriptor
	Description       string
	Path              string
}

type legacyV0ContentTypeDescriptor struct {
	ContentTypeID    int
	Name             string
	Description      string
	DataType         *legacyV0DataType
	ValueDescriptors []*legacyV0ValueDescriptor
}

type legacyV0ChunkStruct struct {
	MetadataOffset  int64
	BodyStartOffset int64
	ChunkSize       int64
}

type legacyV0ChunkMetadata struct {
	ChunkStruct        *legacyV0ChunkStruct
	Producers          []*legacyV0ProducerDescriptor
	StartTimeUnixNano  int64
	EndTimeUnixNano    int64
	StartTicks         int64
	TicksPerNano       float64
	PreviousCheckPoint int64
	ConstantMapPool    legacyV0ConstantMapPool
}

func (c *legacyV0ChunkMetadata) AsUnixNanoTimeStamp(endTicks int64) int64 {
	return c.StartTimeUnixNano + int64(float64(endTicks-c.StartTicks)/c.TicksPerNano)
}

type legacyV0ValueReader interface {
	ReadValue(r *legacyV0Reader, timestamp int64) (any, error)
}

type legacyV0ConstantRef struct {
	pool  *legacyV0ConstantMap
	index int64
}

func (r *legacyV0ConstantRef) Resolve(timestamp int64) any {
	if r == nil || r.pool == nil {
		return nil
	}
	return r.pool.GetResolved(r.index, timestamp)
}

type legacyV0ConstantEntry struct {
	Value        any
	Timestamp    int64
	ResolveState int
}

type legacyV0ConstantMap struct {
	ContentTypeName string
	ValueReader     legacyV0ValueReader
	KeyType         *legacyV0DataType
	entries         map[int64][]*legacyV0ConstantEntry
	allLoaded       bool
}

type legacyV0ConstantMapPool map[int]*legacyV0ConstantMap

func (p legacyV0ConstantMapPool) GetOrInit(contentTypeID int) *legacyV0ConstantMap {
	if cm, ok := p[contentTypeID]; ok {
		return cm
	}
	cm := &legacyV0ConstantMap{entries: make(map[int64][]*legacyV0ConstantEntry)}
	p[contentTypeID] = cm
	return cm
}

func (c *legacyV0ConstantMap) Init(contentTypeName string, reader legacyV0ValueReader, keyType *legacyV0DataType) {
	c.ContentTypeName = contentTypeName
	c.ValueReader = reader
	c.KeyType = keyType
}

func (c *legacyV0ConstantMap) Put(index int64, value any, timestamp int64) {
	c.entries[index] = append(c.entries[index], &legacyV0ConstantEntry{Value: value, Timestamp: timestamp})
}

func (c *legacyV0ConstantMap) ReadValue(r *legacyV0Reader, timestamp int64) error {
	index, err := readLegacyV0ConstantIndex(r, c.KeyType)
	if err != nil {
		return err
	}
	value, err := c.ValueReader.ReadValue(r, timestamp)
	if err != nil {
		return err
	}
	c.Put(index, value, timestamp)
	return nil
}

func (c *legacyV0ConstantMap) SetLoadDone() error {
	if c.ValueReader == nil {
		return fmt.Errorf("nil legacy v0 constant reader")
	}
	c.allLoaded = true
	for _, entries := range c.entries {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp < entries[j].Timestamp
		})
	}
	return nil
}

func (c *legacyV0ConstantMap) Get(index int64, timestamp int64) any {
	if c.allLoaded {
		return c.GetResolved(index, timestamp)
	}
	return &legacyV0ConstantRef{pool: c, index: index}
}

func (c *legacyV0ConstantMap) TouchAll() {
	for idx := range c.entries {
		_ = c.GetResolved(idx, 0)
	}
}

func (c *legacyV0ConstantMap) GetResolved(index int64, timestamp int64) any {
	entries := c.entries[index]
	if len(entries) == 0 {
		return nil
	}
	if timestamp == 0 {
		return entries[len(entries)-1].resolve(c, entries[len(entries)-1].Timestamp)
	}
	for _, entry := range entries {
		if entry.Timestamp >= timestamp {
			return entry.resolve(c, timestamp)
		}
	}
	return entries[len(entries)-1].resolve(c, timestamp)
}

func (e *legacyV0ConstantEntry) resolve(pool *legacyV0ConstantMap, timestamp int64) any {
	switch e.ResolveState {
	case 1:
		return nil
	case 2:
		return e.Value
	}
	e.ResolveState = 1
	resolved := resolveLegacyV0Value(e.Value, timestamp)
	e.Value = convertLegacyV0TypedValue(pool.ContentTypeName, resolved)
	e.ResolveState = 2
	return e.Value
}

func resolveLegacyV0Value(value any, timestamp int64) any {
	switch v := value.(type) {
	case *legacyV0ConstantRef:
		return resolveLegacyV0Value(v.Resolve(timestamp), timestamp)
	case []any:
		values := make([]any, 0, len(v))
		for _, item := range v {
			values = append(values, resolveLegacyV0Value(item, timestamp))
		}
		return values
	default:
		return v
	}
}

type legacyV0PrimitiveReader struct {
	dataType *legacyV0DataType
}

func (r *legacyV0PrimitiveReader) ReadValue(rd *legacyV0Reader, _ int64) (any, error) {
	switch r.dataType {
	case legacyV0Boolean:
		return rd.ReadBool()
	case legacyV0Byte, legacyV0U1:
		return rd.ReadInt8()
	case legacyV0Short, legacyV0U2:
		return rd.ReadShort()
	case legacyV0Integer, legacyV0U4:
		return rd.ReadInt32()
	case legacyV0Long, legacyV0U8:
		return rd.ReadLong()
	case legacyV0Float:
		return rd.ReadFloat32()
	case legacyV0Double:
		return rd.ReadDouble()
	case legacyV0UTF8:
		return rd.ReadString()
	case legacyV0String:
		return rd.ReadUTF16String()
	default:
		return nil, fmt.Errorf("unsupported legacy primitive type [%s]", r.dataType.Name)
	}
}

type legacyV0ConstantReader struct {
	dataType *legacyV0DataType
	pool     *legacyV0ConstantMap
}

func (r *legacyV0ConstantReader) ReadValue(rd *legacyV0Reader, timestamp int64) (any, error) {
	index, err := readLegacyV0ConstantIndex(rd, r.dataType)
	if err != nil {
		return nil, err
	}
	return r.pool.Get(index, timestamp), nil
}

type legacyV0CompositeReader struct {
	readers []legacyV0ValueReader
}

func (r *legacyV0CompositeReader) ReadValue(rd *legacyV0Reader, timestamp int64) (any, error) {
	values := make([]any, 0, len(r.readers))
	for _, reader := range r.readers {
		value, err := reader.ReadValue(rd, timestamp)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

type legacyV0ArrayReader struct {
	reader legacyV0ValueReader
}

func (r *legacyV0ArrayReader) ReadValue(rd *legacyV0Reader, timestamp int64) (any, error) {
	size, err := rd.ReadInt32()
	if err != nil {
		return nil, err
	}
	if size < 0 {
		return nil, fmt.Errorf("illegal legacy array size %d", size)
	}
	values := make([]any, 0, size)
	for i := int32(0); i < size; i++ {
		value, err := r.reader.ReadValue(rd, timestamp)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func readLegacyV0ConstantIndex(r *legacyV0Reader, dt *legacyV0DataType) (int64, error) {
	switch dt {
	case legacyV0Byte, legacyV0U1:
		n, err := r.ReadInt8()
		return int64(n), err
	case legacyV0Short, legacyV0U2:
		n, err := r.ReadShort()
		return int64(n), err
	case legacyV0Integer, legacyV0U4, legacyV0Float:
		n, err := r.ReadInt32()
		return int64(n), err
	case legacyV0Long, legacyV0U8, legacyV0Double:
		return r.ReadLong()
	default:
		return 0, fmt.Errorf("unsupported legacy constant index type [%s]", dt.Name)
	}
}

func createLegacyV0PrimitiveReader(dt *legacyV0DataType, contentType int, metadata *legacyV0ChunkMetadata) (legacyV0ValueReader, error) {
	if dt.Numeric {
		switch contentType {
		case legacyV0ContentTypeEpochMillis, legacyV0ContentTypeMillis, legacyV0ContentTypeNanos,
			legacyV0ContentTypeTicks, legacyV0ContentTypeAddress, legacyV0ContentTypeMemory,
			legacyV0ContentTypePoolNone:
			return &legacyV0PrimitiveReader{dataType: dt}, nil
		default:
			return &legacyV0ConstantReader{dataType: dt, pool: metadata.ConstantMapPool.GetOrInit(contentType)}, nil
		}
	}
	return &legacyV0PrimitiveReader{dataType: dt}, nil
}

func createLegacyV0ReaderFromValueDescriptors(vds []*legacyV0ValueDescriptor, metadata *legacyV0ChunkMetadata) (legacyV0ValueReader, error) {
	if len(vds) == 1 {
		return createLegacyV0Reader(vds[0], metadata)
	}
	readers := make([]legacyV0ValueReader, 0, len(vds))
	for _, vd := range vds {
		reader, err := createLegacyV0Reader(vd, metadata)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return &legacyV0CompositeReader{readers: readers}, nil
}

func createLegacyV0Reader(vd *legacyV0ValueDescriptor, metadata *legacyV0ChunkMetadata) (legacyV0ValueReader, error) {
	switch {
	case vd.DataType.IsPrimitive():
		return createLegacyV0PrimitiveReader(vd.DataType, vd.ContentType, metadata)
	case vd.DataType == legacyV0Array:
		elemType, err := vd.ArrayElemDatatype()
		if err != nil {
			return nil, err
		}
		reader, err := createLegacyV0PrimitiveReader(elemType, vd.ContentType, metadata)
		if err != nil {
			return nil, err
		}
		return &legacyV0ArrayReader{reader: reader}, nil
	case vd.DataType == legacyV0StructArray:
		children, err := vd.Children()
		if err != nil {
			return nil, err
		}
		reader, err := createLegacyV0ReaderFromValueDescriptors(children, metadata)
		if err != nil {
			return nil, err
		}
		return &legacyV0ArrayReader{reader: reader}, nil
	case vd.DataType == legacyV0Struct:
		children, err := vd.Children()
		if err != nil {
			return nil, err
		}
		return createLegacyV0ReaderFromValueDescriptors(children, metadata)
	default:
		return nil, fmt.Errorf("unsupported legacy descriptor type [%s]", vd.DataType.Name)
	}
}

type legacyV0FieldSpec struct {
	Name      string
	ClassName string
}

type legacyV0EventTypeEntry struct {
	TypeName     string
	HasStartTime bool
	Parsers      []legacyV0ValueReader
}

type legacyV0ClassRegistry struct {
	nextID   int64
	classes  map[string]*ClassMetadata
	classMap ClassMap
}

func newLegacyV0ClassRegistry() *legacyV0ClassRegistry {
	return &legacyV0ClassRegistry{
		nextID:   -1,
		classes:  make(map[string]*ClassMetadata),
		classMap: make(ClassMap),
	}
}

func (r *legacyV0ClassRegistry) ensureClass(name, superType string) *ClassMetadata {
	if classMeta, ok := r.classes[name]; ok {
		if classMeta.SuperType == "" && superType != "" {
			classMeta.SuperType = superType
		}
		return classMeta
	}
	classMeta := &ClassMetadata{ID: r.nextID, Name: name, SuperType: superType}
	r.nextID--
	r.classes[name] = classMeta
	r.classMap[int(classMeta.ID)] = classMeta
	return classMeta
}

func (r *legacyV0ClassRegistry) finalize() ClassMap {
	return r.classMap
}

func buildLegacyV0EventClassMetadata(registry *legacyV0ClassRegistry, eventType string, fields []legacyV0FieldSpec) *ClassMetadata {
	classMeta := registry.ensureClass(eventType, EventSuperType)
	classMeta.Fields = make([]FieldMetadata, 0, len(fields))
	for _, field := range fields {
		fieldClass := registry.ensureClass(field.ClassName, "")
		classMeta.Fields = append(classMeta.Fields, FieldMetadata{
			Class: fieldClass.ID,
			Name:  field.Name,
		})
	}
	return classMeta
}

func buildLegacyV0FieldSpecs(etd *legacyV0EventTypeDescriptor, contentTypes map[int]*legacyV0ContentTypeDescriptor) ([]legacyV0FieldSpec, error) {
	fields := make([]legacyV0FieldSpec, 0, len(etd.ValueDescriptors)+4)
	fields = append(fields, legacyV0FieldSpec{Name: "startTime", ClassName: "long"})
	if etd.HasStartTime {
		fields = append(fields, legacyV0FieldSpec{Name: "duration", ClassName: "long"})
	}
	if etd.HasThread {
		fields = append(fields, legacyV0FieldSpec{Name: "eventThread", ClassName: "java.lang.Thread"})
	}
	if etd.CanHaveStacktrace {
		fields = append(fields, legacyV0FieldSpec{Name: "stackTrace", ClassName: "jdk.types.StackTrace"})
	}
	payload, err := flattenLegacyV0Fields(etd.ValueDescriptors, contentTypes, "")
	if err != nil {
		return nil, err
	}
	fields = append(fields, payload...)
	return fields, nil
}

func flattenLegacyV0Fields(vds []*legacyV0ValueDescriptor, contentTypes map[int]*legacyV0ContentTypeDescriptor, prefix string) ([]legacyV0FieldSpec, error) {
	fields := make([]legacyV0FieldSpec, 0, len(vds))
	for _, descriptor := range vds {
		name := descriptor.Identifier
		if prefix != "" {
			name = prefix + "." + descriptor.Identifier
		}
		if descriptor.DataType == legacyV0Struct {
			children, err := descriptor.Children()
			if err != nil {
				return nil, err
			}
			childFields, err := flattenLegacyV0Fields(children, contentTypes, name)
			if err != nil {
				return nil, err
			}
			fields = append(fields, childFields...)
			continue
		}
		fields = append(fields, legacyV0FieldSpec{
			Name:      name,
			ClassName: legacyV0FieldClass(descriptor, contentTypes),
		})
	}
	return fields, nil
}

func legacyV0FieldClass(vd *legacyV0ValueDescriptor, contentTypes map[int]*legacyV0ContentTypeDescriptor) string {
	if contentType, ok := contentTypes[vd.ContentType]; ok {
		if className := legacyV0ContentTypeClassName(contentType.Name); className != "" {
			return className
		}
	}
	switch vd.DataType {
	case legacyV0Boolean:
		return "boolean"
	case legacyV0Byte, legacyV0U1:
		return "byte"
	case legacyV0Short, legacyV0U2:
		return "short"
	case legacyV0Integer, legacyV0U4:
		return "int"
	case legacyV0Long, legacyV0U8:
		return "long"
	case legacyV0Float:
		return "float"
	case legacyV0Double:
		return "double"
	case legacyV0String, legacyV0UTF8:
		return "java.lang.String"
	default:
		return "java.lang.Object"
	}
}

func legacyV0ContentTypeClassName(name string) string {
	switch name {
	case "Thread", "JavaThread", "VMThread":
		return "java.lang.Thread"
	case "StackTrace":
		return "jdk.types.StackTrace"
	case "Class":
		return "java.lang.Class"
	case "Method":
		return "jdk.types.Method"
	case "Package":
		return "jdk.types.Package"
	case "ThreadState":
		return "jdk.types.ThreadState"
	case "Symbol", "UTFConstant", "GCName", "GCCause", "GCWhen":
		return "java.lang.String"
	default:
		return ""
	}
}

func convertLegacyV0TypedValue(contentTypeName string, value any) any {
	if value == nil {
		return nil
	}
	switch contentTypeName {
	case "Thread", "VMThread":
		switch v := value.(type) {
		case string:
			return &Thread{OsName: v, JavaName: v}
		case []any:
			if len(v) > 0 {
				if name, ok := v[0].(string); ok {
					return &Thread{OsName: name, JavaName: name}
				}
			}
		}
	case "JavaThread":
		if values, ok := value.([]any); ok && len(values) > 0 {
			thread, _ := values[0].(*Thread)
			if thread != nil {
				if len(values) > 1 {
					if group, ok := values[1].(string); ok && group != "" && thread.JavaName == "" {
						thread.JavaName = group
					}
				}
				return thread
			}
		}
	case "ThreadState":
		if name, ok := value.(string); ok {
			return &ThreadState{Name: name}
		}
	case "Package":
		if name, ok := value.(string); ok {
			return &Package{Name: &Symbol{String: normalizeLegacyV0JavaName(name)}}
		}
	case "Class":
		if values, ok := value.([]any); ok {
			class := &Class{}
			if len(values) > 1 {
				if name, ok := values[1].(string); ok && name != "" {
					class.Name = &Symbol{String: normalizeLegacyV0JavaName(name)}
				}
			}
			if len(values) > 2 {
				if pkg, ok := values[2].(*Package); ok {
					class.Package = pkg
				}
			}
			if len(values) > 3 {
				switch v := values[3].(type) {
				case int32:
					class.Modifiers = int64(v)
				case int64:
					class.Modifiers = v
				}
			}
			if class.Name != nil {
				return class
			}
		}
	case "Method":
		if values, ok := value.([]any); ok {
			method := &Method{}
			if len(values) > 0 {
				if class, ok := values[0].(*Class); ok {
					method.Type = class
				}
			}
			if len(values) > 1 {
				if name, ok := values[1].(string); ok && name != "" {
					method.Name = &Symbol{String: name}
				}
			}
			if len(values) > 2 {
				if descriptor, ok := values[2].(string); ok && descriptor != "" {
					method.Descriptor = &Symbol{String: descriptor}
				}
			}
			if len(values) > 3 {
				switch v := values[3].(type) {
				case int32:
					method.Modifiers = v
				case int64:
					method.Modifiers = int32(v)
				}
			}
			if method.Name != nil || method.Type != nil {
				return method
			}
		}
	case "StackTrace":
		if values, ok := value.([]any); ok {
			stack := &StackTrace{}
			frameIndex := -1
			truncatedIndex := -1
			if len(values) >= 2 {
				for i, candidate := range values {
					switch candidate.(type) {
					case bool:
						if truncatedIndex < 0 {
							truncatedIndex = i
						}
					case []any:
						if frameIndex < 0 {
							frameIndex = i
						}
					}
				}
			} else if len(values) == 1 {
				if _, ok := values[0].([]any); ok {
					frameIndex = 0
				}
			}
			if truncatedIndex >= 0 {
				stack.Truncated, _ = values[truncatedIndex].(bool)
			}
			if frameIndex >= 0 {
				if frames, ok := values[frameIndex].([]any); ok {
					stack.Frames = make([]*StackFrame, 0, len(frames))
					for _, frameValue := range frames {
						frameValues, ok := frameValue.([]any)
						if !ok {
							continue
						}
						frame := &StackFrame{}
						if len(frameValues) > 0 {
							if method, ok := frameValues[0].(*Method); ok {
								frame.Method = method
							}
						}
						if len(frameValues) > 1 {
							switch v := frameValues[1].(type) {
							case int32:
								frame.LineNumber = v
							case int64:
								frame.LineNumber = int32(v)
							}
						}
						stack.Frames = append(stack.Frames, frame)
					}
				}
			}
			return stack
		}
	}
	return value
}

func normalizeLegacyV0Scalar(value any) any {
	switch v := value.(type) {
	case *Thread, *ThreadState, *Class, *Method, *Package, *StackTrace, *Symbol:
		return v
	case string, bool, int8, int16, int32, int64, float32, float64:
		return v
	default:
		return value
	}
}

func wrapLegacyV0Value(className string, value any) ParseResolvable {
	value = normalizeLegacyV0Scalar(value)
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case *Thread:
		return v
	case *ThreadState:
		return v
	case *Class:
		return v
	case *Method:
		return v
	case *Package:
		return v
	case *StackTrace:
		return v
	case *Symbol:
		return v
	case string:
		switch className {
		case "java.lang.String":
			s := String(v)
			return &s
		case "jdk.types.ThreadState":
			return &ThreadState{Name: v}
		}
	case bool:
		b := Boolean(v)
		return &b
	case int8:
		n := Byte(v)
		return &n
	case int16:
		n := Short(v)
		return &n
	case int32:
		n := Int(v)
		return &n
	case int64:
		n := Long(v)
		return &n
	}
	return WrapRawValue(value)
}

func readLegacyV0Relations(metadataInput *legacyV0Reader) ([]string, error) {
	size, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, err
	}
	relations := make([]string, 0, size)
	for i := int32(0); i < size; i++ {
		relation, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func readLegacyV0DataStructures(metadataInput *legacyV0Reader, relations []string) (*legacyV0DataStructureList, error) {
	size, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, err
	}
	list := make(legacyV0DataStructureList, 0, size)
	for idx := int32(0); idx < size; idx++ {
		fieldCount, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		fields := make([]*legacyV0ValueDescriptor, 0, fieldCount)
		for i := int32(0); i < fieldCount; i++ {
			id, err := metadataInput.ReadString()
			if err != nil {
				return nil, err
			}
			name, err := metadataInput.ReadString()
			if err != nil {
				return nil, err
			}
			description, err := metadataInput.ReadString()
			if err != nil {
				return nil, err
			}
			transition, err := metadataInput.ReadByte()
			if err != nil {
				return nil, err
			}
			dtIdx, err := metadataInput.ReadByte()
			if err != nil {
				return nil, err
			}
			contentType, err := metadataInput.ReadInt32()
			if err != nil {
				return nil, err
			}
			structureIndex, err := metadataInput.ReadInt32()
			if err != nil {
				return nil, err
			}
			relation := ""
			dataType := legacyV0DataTypes[dtIdx]
			if dataType.IsPrimitive() && structureIndex > 0 {
				relation = relations[structureIndex-1]
			}
			if _, err := metadataInput.Skip(legacyV0IntegerSize); err != nil {
				return nil, err
			}
			fields = append(fields, &legacyV0ValueDescriptor{
				Identifier:     id,
				Name:           name,
				Description:    description,
				Transition:     int(transition),
				DataType:       dataType,
				RelationalKey:  relation,
				ContentType:    int(contentType),
				DataStructures: &list,
				StructureIndex: int(structureIndex),
			})
		}
		list.Add(&legacyV0DataStructure{ValueDescriptors: fields})
	}
	return &list, nil
}

func readLegacyV0EventTypes(metadataInput *legacyV0Reader, structures *legacyV0DataStructureList) ([]*legacyV0EventTypeDescriptor, error) {
	size, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, err
	}
	eventTypes := make([]*legacyV0EventTypeDescriptor, 0, size)
	for i := int32(0); i < size; i++ {
		id, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		label, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		path, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		hasStartTime, err := metadataInput.ReadBool()
		if err != nil {
			return nil, err
		}
		hasThread, err := metadataInput.ReadBool()
		if err != nil {
			return nil, err
		}
		canHaveStackTrace, err := metadataInput.ReadBool()
		if err != nil {
			return nil, err
		}
		requestable, err := metadataInput.ReadBool()
		if err != nil {
			return nil, err
		}
		structureIndex, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		if _, err := metadataInput.Skip(legacyV0IntegerSize); err != nil {
			return nil, err
		}
		eventTypes = append(eventTypes, &legacyV0EventTypeDescriptor{
			Identifier:        int(id),
			Label:             label,
			HasStartTime:      hasStartTime,
			HasThread:         hasThread,
			CanHaveStacktrace: canHaveStackTrace,
			IsRequestable:     requestable,
			ValueDescriptors:  structures.Get(int64(structureIndex)).ValueDescriptors,
			Description:       description,
			Path:              path,
		})
	}
	return eventTypes, nil
}

func readLegacyV0ContentTypes(metadataInput *legacyV0Reader, structures *legacyV0DataStructureList) ([]*legacyV0ContentTypeDescriptor, error) {
	size, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, err
	}
	contentTypes := make([]*legacyV0ContentTypeDescriptor, 0, size)
	for i := int32(0); i < size; i++ {
		id, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		name, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		dtIdx, err := metadataInput.ReadByte()
		if err != nil {
			return nil, err
		}
		structureIndex, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		contentTypes = append(contentTypes, &legacyV0ContentTypeDescriptor{
			ContentTypeID:    int(id),
			Name:             name,
			Description:      description,
			DataType:         legacyV0DataTypes[dtIdx],
			ValueDescriptors: structures.Get(int64(structureIndex)).ValueDescriptors,
		})
	}
	return contentTypes, nil
}

func readLegacyV0Producers(metadataInput *legacyV0Reader) ([]*legacyV0ProducerDescriptor, error) {
	size, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, err
	}
	producers := make([]*legacyV0ProducerDescriptor, 0, size)
	for i := int32(0); i < size; i++ {
		id, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, err
		}
		name, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		uri, err := metadataInput.ReadString()
		if err != nil {
			return nil, err
		}
		relations, err := readLegacyV0Relations(metadataInput)
		if err != nil {
			return nil, err
		}
		structures, err := readLegacyV0DataStructures(metadataInput, relations)
		if err != nil {
			return nil, err
		}
		eventTypes, err := readLegacyV0EventTypes(metadataInput, structures)
		if err != nil {
			return nil, err
		}
		contentTypes, err := readLegacyV0ContentTypes(metadataInput, structures)
		if err != nil {
			return nil, err
		}
		producers = append(producers, &legacyV0ProducerDescriptor{
			Name:                   name,
			Description:            description,
			URI:                    uri,
			Identifier:             int(id),
			EventTypeDescriptors:   eventTypes,
			ContentTypeDescriptors: contentTypes,
		})
	}
	return producers, nil
}

func parseLegacyV0ChunkStruct(chunkInput *legacyV0Reader) (*legacyV0ChunkStruct, error) {
	metadataOffset, err := chunkInput.ReadLong()
	if err != nil {
		return nil, err
	}
	bodyStartOffset := chunkInput.Pos()
	lengthBytes := make([]byte, legacyV0IntegerSize)
	if _, err := chunkInput.ReadAt(lengthBytes, metadataOffset); err != nil {
		return nil, err
	}
	metadataLength := int64(binary.BigEndian.Uint32(lengthBytes))
	return &legacyV0ChunkStruct{
		MetadataOffset:  metadataOffset,
		BodyStartOffset: int64(bodyStartOffset),
		ChunkSize:       metadataOffset + metadataLength,
	}, nil
}

func parseLegacyV0ChunkMetadata(metadataInput *legacyV0Reader) (*legacyV0ChunkMetadata, error) {
	if _, err := metadataInput.Skip(legacyV0IntegerSize * 2); err != nil {
		return nil, err
	}
	producers, err := readLegacyV0Producers(metadataInput)
	if err != nil {
		return nil, err
	}
	startTimeMillis, err := metadataInput.ReadLong()
	if err != nil {
		return nil, err
	}
	endTimeMillis, err := metadataInput.ReadLong()
	if err != nil {
		return nil, err
	}
	startTicks, err := metadataInput.ReadLong()
	if err != nil {
		return nil, err
	}
	ticksPerSecond, err := metadataInput.ReadLong()
	if err != nil {
		return nil, err
	}
	prevCheckpoint, err := metadataInput.ReadLong()
	if err != nil {
		return nil, err
	}
	return &legacyV0ChunkMetadata{
		Producers:          producers,
		StartTimeUnixNano:  startTimeMillis * 1e6,
		EndTimeUnixNano:    endTimeMillis * 1e6,
		StartTicks:         startTicks,
		TicksPerNano:       float64(ticksPerSecond) / 1e9,
		PreviousCheckPoint: prevCheckpoint,
		ConstantMapPool:    make(legacyV0ConstantMapPool),
	}, nil
}

func getLegacyV0EventReaders(typeName string, descriptor *legacyV0EventTypeDescriptor, metadata *legacyV0ChunkMetadata) ([]legacyV0ValueReader, error) {
	_ = typeName
	readers := make([]legacyV0ValueReader, 0, len(descriptor.ValueDescriptors)+2)
	if descriptor.HasThread {
		readers = append(readers, &legacyV0ConstantReader{dataType: legacyV0U4, pool: metadata.ConstantMapPool.GetOrInit(legacyV0ContentTypeThreadID)})
	}
	if descriptor.CanHaveStacktrace {
		readers = append(readers, &legacyV0ConstantReader{dataType: legacyV0U8, pool: metadata.ConstantMapPool.GetOrInit(legacyV0ContentTypeStackTrace)})
	}
	for _, field := range descriptor.ValueDescriptors {
		if field.DataType == legacyV0Struct {
			children, err := field.Children()
			if err != nil {
				return nil, err
			}
			for _, child := range children {
				reader, err := createLegacyV0Reader(child, metadata)
				if err != nil {
					return nil, err
				}
				readers = append(readers, reader)
			}
			continue
		}
		reader, err := createLegacyV0Reader(field, metadata)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

func appendLegacyV0Field(ue *UnsupportedEvent, name string, value ParseResolvable) {
	if value == nil {
		return
	}
	if ue.Fields == nil {
		ue.Fields = make(map[string][]ParseResolvable)
	}
	ue.Fields[name] = append(ue.Fields[name], value)
}

func (c *Chunk) parseLegacyV0(version Version, r io.Reader, _ *ChunkParseOptions) error {
	_ = version
	chunkInput := newLegacyV0Reader(r)
	if _, err := chunkInput.FillTo(8); err != nil {
		return fmt.Errorf("unable to read legacy v0 header: %w", err)
	}
	if _, err := chunkInput.Skip(8); err != nil {
		return fmt.Errorf("unable to advance legacy v0 reader beyond magic/version: %w", err)
	}
	chunkStruct, err := parseLegacyV0ChunkStruct(chunkInput)
	if err != nil {
		return fmt.Errorf("unable to parse legacy v0 chunk struct: %w", err)
	}
	if _, err := chunkInput.FillTo(int(chunkStruct.ChunkSize)); err != nil {
		return fmt.Errorf("unable to read legacy v0 chunk data: %w", err)
	}
	metadataInput := newLegacyV0Reader(bytes.NewReader(chunkInput.Buf()[chunkStruct.MetadataOffset:chunkStruct.ChunkSize]))
	chunkMetadata, err := parseLegacyV0ChunkMetadata(metadataInput)
	if err != nil {
		return fmt.Errorf("unable to parse legacy v0 metadata: %w", err)
	}
	chunkMetadata.ChunkStruct = chunkStruct

	c.Header = Header{
		ChunkSize:      chunkStruct.ChunkSize,
		MetadataOffset: chunkStruct.MetadataOffset,
		StartTimeNanos: chunkMetadata.StartTimeUnixNano,
		DurationNanos:  chunkMetadata.EndTimeUnixNano - chunkMetadata.StartTimeUnixNano,
		StartTicks:     chunkMetadata.StartTicks,
		TicksPerSecond: int64(chunkMetadata.TicksPerNano * 1e9),
	}
	c.Metadata = MetadataEvent{
		StartTime: chunkMetadata.StartTimeUnixNano,
		Duration:  chunkMetadata.EndTimeUnixNano - chunkMetadata.StartTimeUnixNano,
	}

	constantTypes := make(map[int64]string)
	for _, producer := range chunkMetadata.Producers {
		for _, descriptor := range producer.EventTypeDescriptors {
			typeName := normalizeLegacyV0EventName(producer.URI + descriptor.Path)
			constantTypes[int64(descriptor.Identifier)] = typeName
		}
		for _, contentType := range producer.ContentTypeDescriptors {
			reader, err := createLegacyV0ReaderFromValueDescriptors(contentType.ValueDescriptors, chunkMetadata)
			if err != nil {
				return fmt.Errorf("unable to create legacy v0 constant reader for [%s]: %w", contentType.Name, err)
			}
			chunkMetadata.ConstantMapPool.GetOrInit(contentType.ContentTypeID).Init(contentType.Name, reader, contentType.DataType)
		}
	}

	prevCheckpoint := chunkMetadata.PreviousCheckPoint
	for prevCheckpoint != 0 {
		lengthBytes := make([]byte, legacyV0IntegerSize)
		if _, err := chunkInput.ReadAt(lengthBytes, prevCheckpoint); err != nil {
			return fmt.Errorf("unable to read legacy checkpoint size: %w", err)
		}
		checkpointSize := binary.BigEndian.Uint32(lengthBytes)
		checkpointEnd := prevCheckpoint + int64(checkpointSize)
		checkpointReader := newLegacyV0Reader(bytes.NewReader(chunkInput.Buf()[prevCheckpoint+(legacyV0IntegerSize*2) : checkpointEnd]))
		ticks, err := checkpointReader.ReadLong()
		if err != nil {
			return fmt.Errorf("unable to read legacy checkpoint ticks: %w", err)
		}
		timestamp := chunkMetadata.AsUnixNanoTimeStamp(ticks)
		prevCheckpoint, err = checkpointReader.ReadLong()
		if err != nil {
			return fmt.Errorf("unable to read legacy next checkpoint offset: %w", err)
		}
		for {
			contentTypeID, err := checkpointReader.ReadInt32()
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("unable to read legacy content type id: %w", err)
			}
			count, err := checkpointReader.ReadInt32()
			if err != nil {
				return fmt.Errorf("unable to read legacy checkpoint entry count: %w", err)
			}
			cm := chunkMetadata.ConstantMapPool.GetOrInit(int(contentTypeID))
			for i := int32(0); i < count; i++ {
				if err := cm.ReadValue(checkpointReader, timestamp); err != nil {
					return fmt.Errorf("unable to read legacy constant value: %w", err)
				}
			}
		}
	}

	for _, cm := range chunkMetadata.ConstantMapPool {
		if err := cm.SetLoadDone(); err != nil {
			return err
		}
	}
	if threadPool := chunkMetadata.ConstantMapPool[legacyV0ContentTypeJavaThread]; threadPool != nil {
		threadPool.TouchAll()
	}

	registry := newLegacyV0ClassRegistry()
	eventTypes := make(map[int]*legacyV0EventTypeEntry)
	eventClasses := make(map[int]*ClassMetadata)
	eventFields := make(map[int][]legacyV0FieldSpec)

	for _, producer := range chunkMetadata.Producers {
		contentTypes := make(map[int]*legacyV0ContentTypeDescriptor, len(producer.ContentTypeDescriptors))
		for _, ct := range producer.ContentTypeDescriptors {
			contentTypes[ct.ContentTypeID] = ct
		}
		for _, descriptor := range producer.EventTypeDescriptors {
			typeName := normalizeLegacyV0EventName(producer.URI + descriptor.Path)
			fields, err := buildLegacyV0FieldSpecs(descriptor, contentTypes)
			if err != nil {
				return fmt.Errorf("unable to build legacy v0 field specs for [%s]: %w", typeName, err)
			}
			readers, err := getLegacyV0EventReaders(typeName, descriptor, chunkMetadata)
			if err != nil {
				return fmt.Errorf("unable to create legacy v0 event readers for [%s]: %w", typeName, err)
			}
			eventTypes[descriptor.Identifier] = &legacyV0EventTypeEntry{
				TypeName:     typeName,
				HasStartTime: descriptor.HasStartTime,
				Parsers:      readers,
			}
			eventClasses[descriptor.Identifier] = buildLegacyV0EventClassMetadata(registry, typeName, fields)
			eventFields[descriptor.Identifier] = fields
		}
	}

	c.Events = make([]EventParseable, 0)

	bodyReader := newLegacyV0Reader(bytes.NewReader(chunkInput.Buf()[chunkStruct.BodyStartOffset:chunkStruct.MetadataOffset]))
	nextPos := bodyReader.Pos()
	for {
		if nextPos > bodyReader.Pos() {
			if _, err := bodyReader.Skip(nextPos - bodyReader.Pos()); err != nil {
				return fmt.Errorf("unable to seek to legacy next event: %w", err)
			}
		}
		eventSize, err := bodyReader.ReadInt32()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("unable to read legacy event size: %w", err)
		}
		nextPos += int(eventSize)
		eventTypeID, err := bodyReader.ReadInt32()
		if err != nil {
			return fmt.Errorf("unable to read legacy event type id: %w", err)
		}
		if eventTypeID == legacyV0MetadataEventTypeIndex || eventTypeID == legacyV0CheckPointEventTypeIndex {
			continue
		}
		entry, ok := eventTypes[int(eventTypeID)]
		if !ok {
			// Best-effort legacy support: ignore unknown event IDs instead of failing the entire file.
			continue
		}
		classMeta := eventClasses[int(eventTypeID)]
		fieldSpecs := eventFields[int(eventTypeID)]

		endTicks, err := bodyReader.ReadLong()
		if err != nil {
			return fmt.Errorf("unable to read legacy event end ticks: %w", err)
		}
		endTime := chunkMetadata.AsUnixNanoTimeStamp(endTicks)
		startTime := endTime
		event := &UnsupportedEvent{}
		event.SetMetadata(classMeta)
		appendLegacyV0Field(event, "startTime", WrapRawValue(startTime))

		fieldOffset := 1
		if entry.HasStartTime {
			startTicks, err := bodyReader.ReadLong()
			if err != nil {
				return fmt.Errorf("unable to read legacy event start ticks: %w", err)
			}
			startTime = chunkMetadata.AsUnixNanoTimeStamp(startTicks)
			event.Fields = make(map[string][]ParseResolvable)
			appendLegacyV0Field(event, "startTime", WrapRawValue(startTime))
			appendLegacyV0Field(event, "duration", WrapRawValue(endTime-startTime))
			fieldOffset++
		}

		for idx, reader := range entry.Parsers {
			value, err := reader.ReadValue(bodyReader, endTime)
			if err != nil {
				return fmt.Errorf("unable to read legacy event value at index [%d]: %w", idx, err)
			}
			fieldIndex := fieldOffset + idx
			if fieldIndex >= len(fieldSpecs) {
				continue
			}
			appendLegacyV0Field(event, fieldSpecs[fieldIndex].Name, wrapLegacyV0Value(fieldSpecs[fieldIndex].ClassName, value))
		}
		c.Events = append(c.Events, event)
	}

	return nil
}

func normalizeLegacyV0EventName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if idx := strings.Index(name, "hotspot/"); idx >= 0 {
		suffix := strings.Trim(name[idx+len("hotspot/"):], "/")
		switch suffix {
		case "jvm/vm/prof/execution_sample":
			return "jdk.ExecutionSample"
		case "jvm/java/thread_park":
			return "jdk.ThreadPark"
		case "jvm/java/thread_sleep":
			return "jdk.ThreadSleep"
		case "jvm/java/monitor_wait":
			return "jdk.JavaMonitorWait"
		case "jfr-info/recordings/recording":
			return "jdk.ActiveRecording"
		case "jfr-info/recordings/recording_setting":
			return "jdk.ActiveSetting"
		}
	}
	return strings.ReplaceAll(name, "/", ".")
}

func normalizeLegacyV0JavaName(name string) string {
	name = strings.ReplaceAll(name, "/", ".")
	return strings.TrimSpace(name)
}

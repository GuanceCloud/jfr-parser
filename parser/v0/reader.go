package v0

import (
	"encoding/binary"
	"fmt"
	types2 "github.com/grafana/jfr-parser/common/types"
	"github.com/grafana/jfr-parser/common/units"
	"github.com/grafana/jfr-parser/parser/utils"
	"io"
	"unicode/utf16"
)

const (
	StringCharElementSize = 2
)

// Content types
const (
	CONTENT_TYPE_POOL_NONE    = 0
	CONTENT_TYPE_MEMORY       = 1
	CONTENT_TYPE_EPOCHMILLIS  = 2
	CONTENT_TYPE_MILLIS       = 3
	CONTENT_TYPE_NANOS        = 4
	CONTENT_TYPE_TICKS        = 5
	CONTENT_TYPE_ADDRESS      = 6
	CONTENT_TYPE_THREADID     = 7
	CONTENT_TYPE_JAVATHREADID = 8
	CONTENT_TYPE_STACKTRACEID = 9
	//	CONTENT_TYPE_CLASSID = 10;
	CONTENT_TYPE_PERCENTAGE = 11
	//	CONTENT_TYPE_VMTHREAD = 30;
	//	CONTENT_TYPE_METHOD = 32;
	//	CONTENT_TYPE_GCWHEN = 38;
)

var notInConstantPoolTypes = map[int]bool{
	CONTENT_TYPE_EPOCHMILLIS: true,
	CONTENT_TYPE_MILLIS:      true,
	CONTENT_TYPE_NANOS:       true,
	CONTENT_TYPE_TICKS:       true,
	CONTENT_TYPE_PERCENTAGE:  true,
	CONTENT_TYPE_MEMORY:      true,
	CONTENT_TYPE_ADDRESS:     true,
	CONTENT_TYPE_POOL_NONE:   true,
}

var (
	defaultBooleanReader     = &BooleanReader{}
	defaultUTF16StringReader = &UTF16StringReader{}
	defaultUTF8StringReader  = &UTF8StringReader{}
)

type ValueReader interface {
	ReadValue(r *utils.JFRReader, timestamp int64) (any, error)
	GetValueType() *ContentTypeDescriptor
}

type NopeValueType struct {
}

func (n *NopeValueType) GetValueType() *ContentTypeDescriptor {
	return nil
}

type TypeIdentifierReader struct {
	datatype              *DataType
	labeledIdentifierPool map[int64]*LabeledIdentifier
	NopeValueType
}

func NewTypeIdentifierReader(dataType *DataType, pool map[int64]*LabeledIdentifier) *TypeIdentifierReader {
	return &TypeIdentifierReader{
		datatype:              dataType,
		labeledIdentifierPool: pool,
	}
}

func (t *TypeIdentifierReader) ReadValue(r *utils.JFRReader, timestampNanos int64) (any, error) {
	index, err := ReadConstantIndex(r, t.datatype)
	if err != nil {
		return nil, fmt.Errorf("unable to read type identifier: %w", err)
	}
	return t.labeledIdentifierPool[index], nil
}

type CompositeReader struct {
	readers []ValueReader
	NopeValueType
}

func NewCompositeReader(readers []ValueReader) *CompositeReader {
	return &CompositeReader{
		readers: readers,
	}
}

func (c *CompositeReader) ReadValue(r *utils.JFRReader, timestampNanos int64) (any, error) {
	values := make([]any, 0, len(c.readers))

	for idx, reader := range c.readers {
		v, err := reader.ReadValue(r, timestampNanos)
		if err != nil {
			return values, fmt.Errorf("unable to read value at index [%d]: %w", idx, err)
		}
		values = append(values, v)
	}
	return values, nil
}

type ArrayReader struct {
	reader ValueReader
	NopeValueType
}

func NewArrayReader(r ValueReader) *ArrayReader {
	return &ArrayReader{
		reader: r,
	}
}

func (a *ArrayReader) ReadValue(r *utils.JFRReader, timestampNanos int64) (any, error) {
	arraySize, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to read array size: %w", err)
	}

	if arraySize < 0 || arraySize > 1e7 {
		return nil, fmt.Errorf("illegal array size: %d", arraySize)
	}

	values := make([]any, 0, arraySize)
	for i := int32(0); i < arraySize; i++ {
		v, err := a.reader.ReadValue(r, timestampNanos)
		if err != nil {
			return values, fmt.Errorf("unable to read array value: %w", err)
		}
		values = append(values, v)
	}
	return values, nil
}

type BooleanReader struct {
	NopeValueType
}

func (b *BooleanReader) ReadValue(r *utils.JFRReader, _ int64) (any, error) {
	return r.ReadBool()
}

type UTF16StringReader struct {
	NopeValueType
}

func (s *UTF16StringReader) ReadValue(r *utils.JFRReader, _ int64) (any, error) {
	charSize, err := r.ReadInt32()
	if err != nil {
		return nil, err
	}
	if charSize < 0 || charSize > 1e8 {
		return nil, fmt.Errorf("illegal string char size: %d", charSize)
	}
	length := charSize * StringCharElementSize

	strData := make([]byte, length)
	if _, err = io.ReadFull(r, strData); err != nil {
		return nil, fmt.Errorf("unable to read string bytes: %w", err)
	}

	u16Chars := make([]uint16, 0, length/2)

	for i := 0; i < len(strData)-1; i += 2 {
		u16Chars = append(u16Chars, binary.BigEndian.Uint16(strData[i:i+2]))
	}

	return string(utf16.Decode(u16Chars)), nil
}

type UTF8StringReader struct {
	NopeValueType
}

func (u *UTF8StringReader) ReadValue(r *utils.JFRReader, _ int64) (any, error) {
	length, err := r.ReadUint16()
	if err != nil {
		return nil, fmt.Errorf("unable to read utf8 string length: %w", err)
	}
	strData := make([]byte, length)
	if _, err := io.ReadFull(r, strData); err != nil {
		return nil, fmt.Errorf("unable to read utf8 string bytes: %w", err)
	}
	return string(strData), nil
}

type ConstantReader struct {
	constantMap *ConstantMap
	dataType    *DataType
	NopeValueType
}

func CreateConstantReader(dataType *DataType, constantMap *ConstantMap) *ConstantReader {
	return &ConstantReader{
		constantMap: constantMap,
		dataType:    dataType,
	}
}

func ReadConstantIndex(r *utils.JFRReader, dataType *DataType) (int64, error) {
	var err error

	switch dataType {
	case U1, Byte:
		var n byte
		n, err = r.ReadByte()
		if err == nil {
			return int64(n), nil
		}
	case Short, U2:
		var n int16
		n, err = r.ReadShort()
		if err == nil {
			return int64(n), nil
		}
	case Integer, U4, Float:
		var n int32
		n, err = r.ReadInt32()
		if err == nil {
			return int64(n), nil
		}
	case Long, U8, Double:
		var n int64
		n, err = r.ReadLong()
		if err == nil {
			return n, nil
		}
	default:
		err = fmt.Errorf("unsupport constant index data type: %s", dataType.Name)
	}
	return 0, err
}

func (c *ConstantReader) ReadValue(r *utils.JFRReader, timestampNanos int64) (any, error) {
	index, err := ReadConstantIndex(r, c.dataType)
	if err != nil {
		return nil, fmt.Errorf("unable to read constant pool index: %w", err)
	}
	return c.constantMap.Get(index, timestampNanos), nil
}

type QuantityReader struct {
	dataType *DataType
	unit     *units.Unit
	NopeValueType
}

func NewQuantityReader(dataType *DataType, u *units.Unit) ValueReader {
	return &QuantityReader{
		dataType: dataType,
		unit:     u,
	}
}

func (q *QuantityReader) ReadValue(r *utils.JFRReader, _ int64) (any, error) {
	switch q.dataType {
	case Byte, U1:
		n, err := r.ReadInt8()
		if err != nil {
			return nil, err
		}
		return q.unit.IntQuantity(int64(n)), nil
	case Short, U2:
		n, err := r.ReadShort()
		if err != nil {
			return nil, err
		}
		return q.unit.IntQuantity(int64(n)), nil
	case Integer, U4:
		n, err := r.ReadInt32()
		if err != nil {
			return nil, err
		}
		return q.unit.IntQuantity(int64(n)), nil
	case Long, U8:
		n, err := r.ReadLong()
		if err != nil {
			return nil, err
		}
		return q.unit.IntQuantity(n), nil
	case Float:
		n, err := r.ReadFloat32()
		if err != nil {
			return nil, err
		}
		return q.unit.FloatQuantity(float64(n)), nil
	case Double:
		n, err := r.ReadDouble()
		if err != nil {
			return nil, err
		}
		return q.unit.FloatQuantity(n), nil
	}
	return nil, fmt.Errorf("unsupport data type: %s", q.dataType.Name)
}

func CreatePrimitiveReader(dataType *DataType, contentType int, valueType string, metadata *ChunkMetadata, labeledIdentifierPool map[int64]*LabeledIdentifier) (ValueReader, error) {
	if dataType.Numeric {
		switch contentType {
		case CONTENT_TYPE_EPOCHMILLIS:
			return NewQuantityReader(dataType, units.UnixMilli), nil
		case CONTENT_TYPE_MILLIS:
			return NewQuantityReader(dataType, units.Millisecond), nil
		case CONTENT_TYPE_NANOS:
			return NewQuantityReader(dataType, units.UnixNano), nil
		case CONTENT_TYPE_TICKS:
			return NewQuantityReader(dataType, metadata.UnitTick), nil
		case CONTENT_TYPE_PERCENTAGE:
			return NewQuantityReader(dataType, units.PercentUnity), nil
		case CONTENT_TYPE_MEMORY:
			return NewQuantityReader(dataType, units.Byte), nil
		case CONTENT_TYPE_ADDRESS:
			return NewQuantityReader(dataType, units.AddressUnity), nil
		case CONTENT_TYPE_POOL_NONE:
			if types2.TypeIdentifierValueInterpretation == valueType {
				return NewTypeIdentifierReader(dataType, labeledIdentifierPool), nil
			} else {
				// TODO: get Unit from valueType
			}
			return NewQuantityReader(dataType, units.NumberUnity), nil
		default:
			return CreateConstantReader(dataType, metadata.ConstantMapPool.GetOrInit(contentType)), nil
		}
	} else if dataType == Boolean {
		return defaultBooleanReader, nil
	} else if dataType == String {
		return defaultUTF16StringReader, nil
	} else if dataType == UTF8 {
		return defaultUTF8StringReader, nil
	}

	return nil, fmt.Errorf("unsupport primitive reader datatype: %s and content type: %d", dataType.Name, contentType)
}

func CreateReaderFromValueDescriptors(vds []*ValueDescriptor, metadata *ChunkMetadata, labeledIdentifierPool map[int64]*LabeledIdentifier) (ValueReader, error) {
	if len(vds) == 1 {
		return CreateReader(vds[0], "", metadata, labeledIdentifierPool)
	}

	readers := make([]ValueReader, 0, len(vds))
	for _, vd := range vds {
		reader, err := CreateReader(vd, "", metadata, labeledIdentifierPool)
		if err != nil {
			return nil, fmt.Errorf("unable to create reader from valueDescriptor: %w", err)
		}
		readers = append(readers, reader)
	}
	return NewCompositeReader(readers), nil
}

func CreateReader(vd *ValueDescriptor, valueType string, metadata *ChunkMetadata, labeledIdentifierPool map[int64]*LabeledIdentifier) (ValueReader, error) {
	switch {
	case vd.DataType.IsPrimitive():
		return CreatePrimitiveReader(vd.DataType, vd.ContentType, valueType, metadata, labeledIdentifierPool)
	case vd.DataType == Array:
		elemType, err := vd.ArrayElemDatatype()
		if err != nil {
			return nil, fmt.Errorf("unable to get array element type: %w", err)
		}
		elemReader, err := CreatePrimitiveReader(elemType, vd.ContentType, valueType, metadata, labeledIdentifierPool)
		if err != nil {
			return nil, fmt.Errorf("unable to create array element reader: %w", err)
		}
		return NewArrayReader(elemReader), nil
	case vd.DataType == StructArray:
		vds, err := vd.Children()
		if err != nil {
			return nil, fmt.Errorf("unable to get children ValueDescriptors: %w", err)
		}
		structReader, err := CreateReaderFromValueDescriptors(vds, metadata, labeledIdentifierPool)
		if err != nil {
			return nil, fmt.Errorf("unable to create struct reader from ValueDescriptor set: %w", err)
		}
		return NewArrayReader(structReader), nil
	case vd.DataType == Struct:
		vds, err := vd.Children()
		if err != nil {
			return nil, fmt.Errorf("unable to get children ValueDescriptors: %w", err)
		}
		return CreateReaderFromValueDescriptors(vds, metadata, labeledIdentifierPool)
	default:
		return nil, fmt.Errorf("unsupport datatype [%s] to create ValueReader", vd.DataType.Name)
	}
}

type EventTypeEntry struct {
	TypeID       string
	Parsers      []ValueReader
	HasStartTime bool
}

func NewEventTypeEntry(typeId string, hasStartTime bool, valueReaders []ValueReader) *EventTypeEntry {
	return &EventTypeEntry{
		TypeID:       typeId,
		Parsers:      valueReaders,
		HasStartTime: hasStartTime,
	}
}

func GetValueReaders(typeId string, etd *EventTypeDescriptor, metadata *ChunkMetadata, labeledIdentifierPool map[int64]*LabeledIdentifier) ([]ValueReader, error) {
	_ = typeId
	readers := make([]ValueReader, 0, len(etd.ValueDescriptors)+2)
	if etd.HasThread {
		readers = append(readers, CreateConstantReader(U4, metadata.ConstantMapPool.GetOrInit(CONTENT_TYPE_THREADID)))
	}
	if etd.CanHaveStacktrace {
		readers = append(readers, CreateConstantReader(U8, metadata.ConstantMapPool.GetOrInit(CONTENT_TYPE_STACKTRACEID)))
	}

	for _, descriptor := range etd.ValueDescriptors {
		if descriptor.DataType == Struct {
			vds, err := descriptor.Children()
			if err != nil {
				return readers, fmt.Errorf("unable to get children valueDescriptors: %w", err)
			}
			for _, vd := range vds {
				identifier := descriptor.Identifier + ":" + vd.Identifier
				name := descriptor.Name + ":" + vd.Name
				_, _ = identifier, name
				// TODO getValueType from identifier and name
				reader, err := CreateReader(vd, "", metadata, labeledIdentifierPool)
				if err != nil {
					return readers, fmt.Errorf("unable to create reader for children valueDescriptor: %w", err)
				}
				readers = append(readers, reader)
			}
		} else {
			reader, err := CreateReader(descriptor, "", metadata, labeledIdentifierPool)
			if err != nil {
				return readers, fmt.Errorf("unable to create reader for valueDescriptor: %w", err)
			}
			readers = append(readers, reader)
		}
	}

	return readers, nil
}

package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	types2 "github.com/grafana/jfr-parser/common/types"
	"github.com/grafana/jfr-parser/common/units"
	"github.com/grafana/jfr-parser/parser/utils"
	v0 "github.com/grafana/jfr-parser/parser/v0"
	"io"
)

const (
	JFRVersion0 = 0 // JDK7 & JDK8
	JFRVersion1 = 1 // JDK9 & JDK10
	JFRVersion2 = 2 // JDK11

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
	Pool     map[int64]ParseResolvable
	resolved bool
}
type ClassMap map[int64]*ClassMetadata
type PoolMap map[int64]*CPool
type ChunkEvents map[string]*EventCollection

func (c ChunkEvents) Apply(filter EventFilter) []*GenericEvent {
	var filtered []*GenericEvent

	for _, collection := range c {
		predicate := filter.GetPredicate(collection.ClassMetadata)
		if IsAlwaysFalse(predicate) {
			continue
		} else if IsAlwaysTrue(predicate) {
			filtered = append(filtered, collection.Events...)
		} else {
			for _, event := range collection.Events {
				if predicate.Test(event) {
					filtered = append(filtered, event)
				}
			}
		}
	}
	return filtered
}

type Chunk struct {
	Header   Header
	Metadata ChunkMetadata
	ChunkEvents
}

type EventCollection struct {
	ClassMetadata *ClassMetadata
	Events        []*GenericEvent
}

func (c *EventCollection) Add(e *GenericEvent) {
	c.Events = append(c.Events, e)
}

type ChunkParseOptions struct {
	CPoolProcessor func(meta *ClassMetadata, cpool *CPool)
}

func (c *Chunk) addEvent(e *GenericEvent) {
	if c.ChunkEvents == nil {
		c.ChunkEvents = make(ChunkEvents)
	}

	classMeta := c.Metadata.ClassMap[e.ClassID]

	ec, ok := c.ChunkEvents[classMeta.Name]
	if !ok {
		ec = &EventCollection{
			ClassMetadata: classMeta,
		}
		c.ChunkEvents[classMeta.Name] = ec
	}
	ec.Add(e)
}

func (c *Chunk) ShowClassMeta(name string) {
	for _, classMeta := range c.Metadata.ClassMap {
		if classMeta.Name == name {
			fmt.Printf("simple type: %t, super type: %s\n",
				classMeta.SimpleType, classMeta.SuperType)

			for _, field := range classMeta.Fields {
				fmt.Printf("field name: %s, field label: %s, field class: %s, field description: %s, field constant pool: %t, field is array: %t, field unsigned: %t, field unit: %+#v\n",
					field.Name, field.Label(c.Metadata.ClassMap), c.Metadata.ClassMap[field.ClassID].Name, field.Description(c.Metadata.ClassMap),
					field.ConstantPool, field.IsArray(), field.Unsigned(c.Metadata.ClassMap), field.Unit(c.Metadata.ClassMap))
			}
			break
		}
	}
}

func readRelations(metadataInput *utils.JFRReader) ([]string, error) {
	relationSize, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to read relations size: %w", err)
	}

	relations := make([]string, 0, relationSize)
	for i := int32(0); i < relationSize; i++ {
		relation, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read next relation: %w", err)
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func readDataStructures(metadataInput *utils.JFRReader, relations []string) (*v0.DataStructureList, error) {
	dataStructureSize, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to read data structure size: %w", err)
	}
	dataStructureList := v0.NewDataStructureList(int64(dataStructureSize))

	for idx := int32(0); idx < dataStructureSize; idx++ {
		valueDescriptorSize, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read value descriptor size: %w", err)
		}

		valueDescriptors := make([]*v0.ValueDescriptor, 0, valueDescriptorSize)
		for i := int32(0); i < valueDescriptorSize; i++ {
			id, err := metadataInput.ReadString()
			if err != nil {
				return nil, fmt.Errorf("unable to read value descriptor id: %w", err)
			}
			name, err := metadataInput.ReadString()
			if err != nil {
				return nil, fmt.Errorf("unable to read value descriptor name: %w", err)
			}
			description, err := metadataInput.ReadString()
			if err != nil {
				return nil, fmt.Errorf("unable to read value descriptpor: %w", err)
			}
			transition, err := metadataInput.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("unable to read transition: %w", err)
			}

			datatypeIdx, err := metadataInput.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("unable to read datatype: %w", err)
			}
			dataType := v0.DataTypeEnum[datatypeIdx]

			contentType, err := metadataInput.ReadInt32()
			if err != nil {
				return nil, fmt.Errorf("unable to read content type: %w", err)
			}
			dataStructIndex, err := metadataInput.ReadInt32()
			if err != nil {
				return nil, fmt.Errorf("unable to read data struct index")
			}
			var relation string
			if dataType.IsPrimitive() && dataStructIndex > 0 {
				relation = relations[dataStructIndex-1]
			}

			// expansion, not used now
			if _, err = metadataInput.Skip(v0.IntegerSize); err != nil {
				return nil, fmt.Errorf("unable to skip value descriptor expansion: %w", err)
			}

			valueDescriptors = append(valueDescriptors, &v0.ValueDescriptor{
				Identifier:     id,
				Name:           name,
				Description:    description,
				Transition:     int(transition),
				DataType:       dataType,
				RelationalKey:  relation,
				ContentType:    int(contentType),
				DataStructures: dataStructureList,
				StructureIndex: int(dataStructIndex),
			})
		}

		dataStructureList.Add(&v0.DataStructure{ValueDescriptors: valueDescriptors})
	}
	return dataStructureList, nil
}

func readEventTypes(metadataInput *utils.JFRReader, dataStructures *v0.DataStructureList) ([]*v0.EventTypeDescriptor, error) {
	eventTypeSize, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to read event type size: %w", err)
	}

	eventTypes := make([]*v0.EventTypeDescriptor, 0, eventTypeSize)

	for i := int32(0); i < eventTypeSize; i++ {
		id, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type id: %w", err)
		}
		name, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type name: %w", err)
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type description: %w", err)
		}
		path, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type path: %w", err)
		}
		hasStartTime, err := metadataInput.ReadBool()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type hasStartTime: %w", err)
		}
		hasThread, err := metadataInput.ReadBool()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type hasThread: %w", err)
		}
		canHaveStacktrace, err := metadataInput.ReadBool()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type canHaveStacktrace: %w", err)
		}
		isRequestAble, err := metadataInput.ReadBool()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type isRequestAble: %w", err)
		}

		dataStructureIdx, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read event type dataStructure index: %w", err)
		}

		dataStructure := dataStructures.Get(int64(dataStructureIdx))

		// expansion, not used now
		if _, err = metadataInput.Skip(v0.IntegerSize); err != nil {
			return nil, fmt.Errorf("unable to skip event type expansion bytes: %w", err)
		}

		eventTypes = append(eventTypes, &v0.EventTypeDescriptor{
			Identifier:        int(id),
			Label:             name,
			HasStartTime:      hasStartTime,
			HasThread:         hasThread,
			CanHaveStacktrace: canHaveStacktrace,
			IsRequestAble:     isRequestAble,
			ValueDescriptors:  dataStructure.ValueDescriptors,
			Description:       description,
			Path:              path,
		})
	}

	return eventTypes, nil
}

func readContentTypes(metadataInput *utils.JFRReader, dataStructures *v0.DataStructureList) ([]*v0.ContentTypeDescriptor, error) {
	contentTypeSize, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to read content type size: %w", err)
	}

	contentTypes := make([]*v0.ContentTypeDescriptor, 0, contentTypeSize)

	for i := int32(0); i < contentTypeSize; i++ {
		constantPoolIndex, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read content type constant pool index: %w", err)
		}
		name, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read content type name: %w", err)
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read content type description: %w", err)
		}

		dataTypeIndex, err := metadataInput.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("unable to read content type data type index: %w", err)
		}

		dataType := v0.DataTypeEnum[dataTypeIndex]

		dataStructureIndex, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read content type dataStructure index: %w", err)
		}

		dataStructure := dataStructures.Get(int64(dataStructureIndex))

		contentTypes = append(contentTypes, &v0.ContentTypeDescriptor{
			ContentTypeId:    int(constantPoolIndex),
			Name:             name,
			Description:      description,
			DataType:         dataType,
			ValueDescriptors: dataStructure.ValueDescriptors,
		})
	}

	return contentTypes, nil
}

func readProducers(metadataInput *utils.JFRReader) ([]*v0.ProducerDescriptor, error) {
	producerSize, err := metadataInput.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("unable to resolve producer size: %w", err)
	}

	producers := make([]*v0.ProducerDescriptor, 0, producerSize)

	for i := int32(0); i < producerSize; i++ {
		id, err := metadataInput.ReadInt32()
		if err != nil {
			return nil, fmt.Errorf("unable to read producer identifier: %w", err)
		}
		name, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read producer name: %w", err)
		}
		description, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read producer description: %w", err)
		}
		uri, err := metadataInput.ReadString()
		if err != nil {
			return nil, fmt.Errorf("unable to read producer uri: %w", err)
		}

		relations, err := readRelations(metadataInput)
		if err != nil {
			return nil, fmt.Errorf("unable to read relations: %w", err)
		}

		dataStructures, err := readDataStructures(metadataInput, relations)
		if err != nil {
			return nil, fmt.Errorf("unable to read data structures: %w", err)
		}

		eventTypes, err := readEventTypes(metadataInput, dataStructures)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve event types: %w", err)
		}

		contentTypes, err := readContentTypes(metadataInput, dataStructures)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve content types: %w", err)
		}

		producers = append(producers, &v0.ProducerDescriptor{
			Name:                   name,
			Description:            description,
			UriString:              uri,
			Identifier:             int(id),
			EventTypeDescriptors:   eventTypes,
			ContentTypeDescriptors: contentTypes,
		})
	}
	return producers, nil
}

func parseChunkStruct(chunkInput *utils.JFRReader) (*v0.ChunkStruct, error) {
	metadataOffSet, err := chunkInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read chunk metadata offset: %w", err)
	}

	bodyStartOffset := chunkInput.Pos()

	metadataLenData := make([]byte, 4)
	if _, err = chunkInput.ReadAt(metadataLenData, metadataOffSet); err != nil {
		return nil, fmt.Errorf("unable to read metadata length bytes: %w", err)
	}
	metadataLength := int64(binary.BigEndian.Uint32(metadataLenData))

	return &v0.ChunkStruct{
		MetadataOffset:  metadataOffSet,
		BodyStartOffset: int64(bodyStartOffset),
		ChunkSize:       metadataOffSet + metadataLength,
	}, nil
}

func parseChunkMetadata(metadataInput *utils.JFRReader) (*v0.ChunkMetadata, error) {
	// skip metadata size and event type data
	if _, err := metadataInput.Skip(v0.IntegerSize * 2); err != nil {
		return nil, fmt.Errorf("unable to skip metadata size bytes: %w", err)
	}

	producers, err := readProducers(metadataInput)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve metadata producers: %w", err)
	}

	startTimeMillis, err := metadataInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata startTime unixmillis: %w", err)
	}
	startTimeNanos := startTimeMillis * 1e6

	endTimeMillis, err := metadataInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata endTime unixmillis: %w", err)
	}
	endTimeNanos := endTimeMillis * 1e6

	startTicks, err := metadataInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata startTicks: %w", err)
	}

	ticksPerSec, err := metadataInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata ticksPerSeconds: %w", err)
	}

	ticksPerNano := float64(ticksPerSec) / 1e9

	previousCheckPoint, err := metadataInput.ReadLong()
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata previousCheckPoint: %w", err)
	}

	unitTick := units.Nanosecond.Derived("tick", units.F64(1/ticksPerNano))

	return &v0.ChunkMetadata{
		Producers:          producers,
		StartTimeUnixNano:  startTimeNanos,
		EndTimeUnixNano:    endTimeNanos,
		StartTicks:         startTicks,
		TicksPerNano:       ticksPerNano,
		PreviousCheckPoint: previousCheckPoint,
		UnitTick:           unitTick,
	}, nil
}

func (c *Chunk) Parse(r io.Reader, options *ChunkParseOptions) (err error) {
	chunkInput := utils.NewDataReader(r)

	buf := make([]byte, len(magic))
	if _, err = io.ReadFull(chunkInput, buf); err != nil {
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

	// read chunk's version
	majorVersion, err := chunkInput.ReadUint16()
	if err != nil {
		return fmt.Errorf("unable to read JFR major version: %w", err)
	}
	minorVersion, err := chunkInput.ReadUint16()
	if err != nil {
		return fmt.Errorf("unable to read JFR minor version: %w", err)
	}

	l.Debugf("JFR version: %d.%d\n", majorVersion, minorVersion)

	if majorVersion == JFRVersion0 {
		chunkStruct, err := parseChunkStruct(chunkInput)
		if err != nil {
			return fmt.Errorf("unable to parse chunk struct: %w", err)
		}

		if _, err = chunkInput.FillTo(int(chunkStruct.ChunkSize)); err != nil {
			return fmt.Errorf("unable to read full chunk data: %w", err)
		}

		metadataInput := utils.NewDataReader(bytes.NewReader(chunkInput.Buf()[chunkStruct.MetadataOffset:chunkStruct.ChunkSize]))
		chunkMetadata, err := parseChunkMetadata(metadataInput)
		if err != nil {
			return fmt.Errorf("unable to parse chunk metadata: %w", err)
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
		c.Metadata = ChunkMetadata{
			StartTime: chunkMetadata.StartTimeUnixNano,
			Duration:  chunkMetadata.EndTimeUnixNano - chunkMetadata.StartTimeUnixNano,
			Header:    &c.Header,
		}

		constantMapPool := make(v0.ConstantMapPool)
		constantTypes := make(map[int64]*v0.LabeledIdentifier)
		chunkMetadata.ConstantMapPool = constantMapPool

		// read constants
		for _, producer := range chunkMetadata.Producers {
			for _, etd := range producer.EventTypeDescriptors {
				constantTypes[int64(etd.Identifier)] = &v0.LabeledIdentifier{
					InterfaceId: producer.UriString + etd.Path,
					ImplId:      int64(etd.Identifier),
					Name:        etd.Label,
					Description: etd.Description,
				}
			}

			for _, ct := range producer.ContentTypeDescriptors {
				valueReader, err := v0.CreateReaderFromValueDescriptors(ct.ValueDescriptors, chunkMetadata, constantTypes)
				if err != nil {
					return fmt.Errorf("unable to create constant reader: %w", err)
				}
				constantMapPool.GetOrInit(ct.ContentTypeId).Init(valueReader, ct.DataType)
			}

		}

		prevCheckPoint := chunkMetadata.PreviousCheckPoint

		for prevCheckPoint != 0 {
			lengthBytes := make([]byte, 4)
			if _, err := chunkInput.ReadAt(lengthBytes, prevCheckPoint); err != nil {
				return fmt.Errorf("unable to read checkpoint size data bytes: %w", err)
			}
			checkPointSize := binary.BigEndian.Uint32(lengthBytes)
			checkPointEnd := prevCheckPoint + int64(checkPointSize)
			checkPointReader := utils.NewDataReader(bytes.NewReader(chunkInput.Buf()[prevCheckPoint+(v0.IntegerSize*2) : checkPointEnd]))
			//if _, err = checkPointReader.Skip(v0.IntegerSize); err != nil {
			//	return fmt.Errorf("unable to skip check point event type: %w", err)
			//}
			timestampTicks, err := checkPointReader.ReadLong()
			if err != nil {
				return fmt.Errorf("unable to read checkpoint timestamp ticks: %w", err)
			}
			timestampNanos := chunkMetadata.AsUnixNanoTimeStamp(timestampTicks)
			prevCheckPoint, err = checkPointReader.ReadLong()
			if err != nil {
				return fmt.Errorf("unable to read next checkpoint offset: %w", err)
			}
			for {
				contentTypeId, err := checkPointReader.ReadInt32()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return fmt.Errorf("unable to read constant content type id: %w", err)
				}
				checkPointCount, err := checkPointReader.ReadInt32()
				if err != nil {
					return fmt.Errorf("unable to read entry count: %w", err)
				}
				cm := constantMapPool.GetOrInit(int(contentTypeId))
				for i := int32(0); i < checkPointCount; i++ {
					if err := cm.ReadValue(checkPointReader, timestampNanos); err != nil {
						return fmt.Errorf("unable to read constant value: %w", err)
					}
				}
			}
		}

		for _, constantMap := range constantMapPool {
			if err := constantMap.SetLoadDone(); err != nil {
				return fmt.Errorf("unable to set constant map load done: %w", err)
			}
		}

		threadPool := constantMapPool[v0.CONTENT_TYPE_JAVATHREADID]
		if threadPool != nil {
			threadPool.TouchAll()
		}

		registry := newV0ClassRegistry()
		eventTypes := make(map[int]*v0.EventTypeEntry)
		eventClasses := make(map[int]*ClassMetadata)
		eventFields := make(map[int][]v0FieldSpec)

		for _, producer := range chunkMetadata.Producers {
			contentTypes := make(map[int]*v0.ContentTypeDescriptor, len(producer.ContentTypeDescriptors))
			for _, contentType := range producer.ContentTypeDescriptors {
				contentTypes[contentType.ContentTypeId] = contentType
			}

			for _, descriptor := range producer.EventTypeDescriptors {
				typeID := types2.ConvertToJDK11Type(producer.UriString + descriptor.Path)
				fields, err := buildV0FieldSpecs(descriptor, contentTypes)
				if err != nil {
					return fmt.Errorf("unable to build v0 field metadata for [%s]: %w", typeID, err)
				}
				readers, err := v0.GetValueReaders(typeID, descriptor, chunkMetadata, constantTypes)
				if err != nil {
					return fmt.Errorf("unable to create v0 readers for [%s]: %w", typeID, err)
				}

				eventTypes[descriptor.Identifier] = v0.NewEventTypeEntry(typeID, descriptor.HasStartTime, readers)
				eventClasses[descriptor.Identifier] = buildV0EventClassMetadata(registry, typeID, fields)
				eventFields[descriptor.Identifier] = fields
			}
		}
		c.Metadata.ClassMap = registry.finalize()
		for _, classMeta := range c.Metadata.ClassMap {
			for _, field := range classMeta.Fields {
				field.ChunkHeader = &c.Header
			}
		}

		chunkBodyInput := utils.NewDataReader(bytes.NewReader(chunkInput.Buf()[chunkStruct.BodyStartOffset:chunkStruct.MetadataOffset]))
		nextPos := chunkBodyInput.Pos()

		for {
			if nextPos > chunkBodyInput.Pos() {
				if _, err := chunkBodyInput.Skip(nextPos - chunkBodyInput.Pos()); err != nil {
					return fmt.Errorf("unable to skip to next event pos: %w", err)
				}
			}
			eventSize, err := chunkBodyInput.ReadInt32()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("unable to read next event size: %w", err)
			}

			nextPos += int(eventSize)

			eventTypeId, err := chunkBodyInput.ReadInt32()
			if err != nil {
				return fmt.Errorf("unable to read next event type id: %w", err)
			}

			if eventTypeId == v0.MetadataEventTypeIndex || eventTypeId == v0.CheckPointEventTypeIndex {
				// ignore metadata event and checkpoint event
			} else {
				// Data event
				ep, ok := eventTypes[int(eventTypeId)]
				if !ok {
					return fmt.Errorf("event type [%d] is not describe in the file", eventTypeId)
				}
				classMeta := eventClasses[int(eventTypeId)]
				fieldSpecs := eventFields[int(eventTypeId)]

				endTimeTicks, err := chunkBodyInput.ReadLong()
				if err != nil {
					return fmt.Errorf("unable to read event endtime ticks: %w", err)
				}
				endTimeUnixNano := chunkMetadata.AsUnixNanoTimeStamp(endTimeTicks)
				startTimeUnixNano := endTimeUnixNano
				ge := NewGenericEvent(classMeta.ID, classMeta)

				if ep.HasStartTime {
					startTimeTicks, err := chunkBodyInput.ReadLong()
					if err != nil {
						return fmt.Errorf("unable to read event start time ticks: %w", err)
					}
					startTimeUnixNano = chunkMetadata.AsUnixNanoTimeStamp(startTimeTicks)
					ge.Attributes["duration"] = WrapRawValue(units.Nanosecond.IntQuantity(endTimeUnixNano - startTimeUnixNano))
				}
				ge.Attributes["startTime"] = WrapRawValue(units.UnixNano.IntQuantity(startTimeUnixNano))

				fieldOffset := 1
				if ep.HasStartTime {
					fieldOffset++
				}
				for idx, parser := range ep.Parsers {
					val, err := parser.ReadValue(chunkBodyInput, endTimeUnixNano)
					if err != nil {
						return fmt.Errorf("unable to read event value at index[%d]: %w", idx, err)
					}
					fieldIndex := fieldOffset + idx
					if fieldIndex >= len(fieldSpecs) {
						return fmt.Errorf("event type [%d] field index out of range", eventTypeId)
					}
					ge.Attributes[fieldSpecs[fieldIndex].Name] = wrapV0Value(fieldSpecs[fieldIndex].ClassName, val)
				}
				c.addEvent(ge)
			}
		}

		return nil

	}

	buf = make([]byte, headerSize)
	if _, err = io.ReadFull(chunkInput, buf); err != nil {
		return fmt.Errorf("unable to read chunk header: %w", err)
	}
	if err = c.Header.Parse(NewReader(bytes.NewReader(buf), false)); err != nil {
		return fmt.Errorf("unable to parse chunk header: %w", err)
	}
	c.Header.ChunkSize -= headerSize + 8
	c.Header.MetadataOffset -= headerSize + 8
	c.Header.ConstantPoolOffset -= headerSize + 8
	useCompression := c.Header.Features&1 == 1
	// TODO: assert c.Header.ChunkSize is small enough
	l.Debugf("chunk size: %+#v", c.Header)
	buf = make([]byte, c.Header.ChunkSize)
	if _, err := io.ReadFull(chunkInput, buf); err != nil {
		return fmt.Errorf("unable to read chunk contents: %w", err)
	}

	br := bytes.NewReader(buf)
	rd := NewReader(br, useCompression)
	pointer := int64(0)
	eventsOffset := make(map[int64]int32)

	// Parse metadata
	if _, err := br.Seek(c.Header.MetadataOffset, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	metadataSize, err := rd.VarInt()
	if err != nil {
		return fmt.Errorf("unable to parse chunk metadata size: %w", err)
	}
	eventsOffset[c.Header.MetadataOffset] = metadataSize

	c.Metadata.Header = &c.Header // for resolving class fields unit ticks

	if err := c.Metadata.Parse(rd); err != nil {
		return fmt.Errorf("unable to parse chunk metadata: %w", err)
	}

	// Parse checkpoint event(s)
	if _, err := br.Seek(c.Header.ConstantPoolOffset, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	constantPoolSize := int32(0)
	cpools := make(PoolMap)
	delta := int64(0)
	cp := new(ConstantPoolEvent)
	for {
		size, err := rd.VarInt()
		if err != nil {
			return fmt.Errorf("unable to parse checkpoint event size: %w", err)
		}
		eventsOffset[c.Header.ConstantPoolOffset+delta] = size
		constantPoolSize += size
		if err := cp.Parse(rd, c.Metadata.ClassMap, cpools); err != nil {
			return fmt.Errorf("unable to parse constant pool event: %w", err)
		}
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
			options.CPoolProcessor(c.Metadata.ClassMap[classID], pool)
		}
	}

	// Second pass over constant pools: resolve constants
	if err = ResolveConstants(c.Metadata.ClassMap, cpools); err != nil {
		return err
	}

	// Parse the rest of events
	if _, err := br.Seek(pointer, io.SeekStart); err != nil {
		return fmt.Errorf("unable to seek reader: %w", err)
	}
	for pointer != c.Header.ChunkSize {
		if size, ok := eventsOffset[pointer]; ok {
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
			eventsOffset[pointer] = size
			ge, err := ParseEvent(rd, c.Metadata.ClassMap, cpools)
			if err != nil {
				return fmt.Errorf("unable to parse event: %w", err)
			}
			if ge != nil {
				c.addEvent(ge)
			}
			pointer += int64(size)
		}
	}
	return nil
}

func ResolveConstants(classes ClassMap, poolMap PoolMap) (err error) {

	for classID, pool := range poolMap {
		if pool.resolved {
			continue
		}
		if classes[classID] != nil && classes[classID].SuperType == EventSuperType {
			continue
		}

		for _, t := range pool.Pool {
			if err = t.Resolve(classes, poolMap); err != nil {
				return fmt.Errorf("unable to resolve constants: %w", err)
			}
		}
		pool.resolved = true
	}

	return nil
}

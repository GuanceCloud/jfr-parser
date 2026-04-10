package v0

import (
	"fmt"
	"github.com/grafana/jfr-parser/parser/utils"
	"sort"
)

type ConstantMapPool map[int]*ConstantMap

func (p ConstantMapPool) GetOrInit(contentTypeId int) *ConstantMap {
	cm, ok := p[contentTypeId]
	if !ok {
		cm = &ConstantMap{
			entryListMap: make(map[int64][]*ConstantEntry),
		}
		p[contentTypeId] = cm
	}
	return cm
}

type ConstantEntry struct {
	Value        any
	Timestamp    int64
	ResolveState int
}

func NewConstantEntry(v any, timestamp int64) *ConstantEntry {
	return &ConstantEntry{
		Value:     v,
		Timestamp: timestamp,
	}
}

func (c *ConstantEntry) GetResolved(index int64, cm *ConstantMap) any {
	switch c.ResolveState {
	case 1:
		return nil
	case 0:
		c.ResolveState = 1
		c.Value = ResolveConstant(c.Value, c.Timestamp, cm)
		// TODO: factory map value to detailed type
		c.ResolveState = 2
		return c.Value
	default:
		return c.Value
	}
}

func ResolveConstant(v any, timestampNanos int64, cm *ConstantMap) any {
	if cf, ok := v.(*ConstantReference); ok {
		return ResolveConstant(cf.Resolve(cm, timestampNanos), timestampNanos, cm)
	} else if values, ok := v.([]any); ok {
		for i, value := range values {
			values[i] = ResolveConstant(value, timestampNanos, cm)
		}
		return values
	}
	return v
}

type ConstantMap struct {
	allConstantsLoaded bool

	entryListMap map[int64][]*ConstantEntry

	ValueReader ValueReader
	KeyType     *DataType
}

func (c *ConstantMap) EntryListMap() map[int64][]*ConstantEntry {
	return c.entryListMap
}

func (c *ConstantMap) Init(reader ValueReader, dataType *DataType) {
	c.ValueReader = reader
	c.KeyType = dataType
}

func (c *ConstantMap) GetResolved(index int64, timestampNanos int64) any {
	entryList := c.entryListMap[index]
	if len(entryList) > 0 {
		for _, entry := range entryList {
			if entry.Timestamp >= timestampNanos {
				return entry.GetResolved(index, c)
			}
		}
	}
	// TODO factory provide empty object
	return nil
}

func (c *ConstantMap) Get(index int64, timestampNanos int64) any {
	if c.allConstantsLoaded {
		return c.GetResolved(index, timestampNanos)
	}
	return &ConstantReference{
		index: index,
	}
}

func (c *ConstantMap) Put(index int64, value any, timestamp int64) {
	c.entryListMap[index] = append(c.entryListMap[index], NewConstantEntry(value, timestamp))
}

func (c *ConstantMap) SetLoadDone() error {
	c.allConstantsLoaded = true
	if c.ValueReader == nil {
		return fmt.Errorf("nil constant pool reader is not allowed")
	}
	for _, entries := range c.entryListMap {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp < entries[j].Timestamp
		})
	}
	return nil
}

func (c *ConstantMap) TouchAll() {
	for idx, entries := range c.entryListMap {
		for _, entry := range entries {
			entry.GetResolved(idx, c)
		}
	}
}

func (c *ConstantMap) ReadValue(reader *utils.JFRReader, timestamp int64) error {
	index, err := ReadConstantIndex(reader, c.KeyType)
	if err != nil {
		return fmt.Errorf("unable to read constant pool index: %w", err)
	}
	v, err := c.ValueReader.ReadValue(reader, timestamp)
	if err != nil {
		return fmt.Errorf("unable to read constant value: %w", err)
	}
	c.Put(index, v, timestamp)
	return nil
}

type ConstantReference struct {
	index int64
}

func (c *ConstantReference) Resolve(cm *ConstantMap, timestampNanos int64) any {
	return cm.GetResolved(c.index, timestampNanos)
}

type LabeledIdentifier struct {
	InterfaceId string
	ImplId      int64
	Name        string
	Description string
}

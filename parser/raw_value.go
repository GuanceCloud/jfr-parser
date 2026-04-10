package parser

import "fmt"

// RawValue is a lightweight ParseResolvable wrapper for compatibility paths
// that already decoded a value without going through the regular metadata-based
// parser pipeline.
type RawValue struct {
	Value any
}

func (r *RawValue) Parse(Reader, ClassMap, PoolMap, *ClassMetadata) error {
	return fmt.Errorf("raw values are not parseable")
}

func (*RawValue) Resolve(ClassMap, PoolMap) error {
	return nil
}

func WrapRawValue(v any) ParseResolvable {
	if v == nil {
		return nil
	}
	if resolvable, ok := v.(ParseResolvable); ok {
		return resolvable
	}
	return &RawValue{Value: v}
}

func UnwrapValue(v any) any {
	if raw, ok := v.(*RawValue); ok {
		return raw.Value
	}
	return v
}

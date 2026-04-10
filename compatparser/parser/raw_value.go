package parser

import "fmt"

// RawValue wraps already-decoded compatibility values so they can flow
// through the existing ParseResolvable-based business logic.
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

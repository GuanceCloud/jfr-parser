package common

type StructContentType struct {
	identifier string
	name       string
	attributes []*Attribute[any]
	accessors  map[Attribute[any]]IMemberAccessor[any]
}

func GetAccessor[T any](s *StructContentType, attr Attribute[T]) IMemberAccessor[any] {
	if accessor, ok := s.accessors[Attribute[any](attr)]; ok {
		return accessor
	}
	return nil
}

type EventItem struct {
	eventType *StructContentType
	fields    map[string]interface{}
}

type IMemberAccessor[T any] interface {
	GetMember(item *EventItem) T
}

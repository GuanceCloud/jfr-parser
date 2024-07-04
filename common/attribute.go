package common

type Attribute[T any] struct {
	identifier            string
	name                  string
	description           string
	contentTypeIdentifier string
}

func Attr[T any](identifier, name, description string, contentType StructContentType) *Attribute[T] {
	return &Attribute[T]{
		identifier:  identifier,
		name:        name,
		description: description,
	}
}

func (a *Attribute[T]) GetAccessor(contentType *StructContentType) IMemberAccessor[T] {
	return GetAccessor(contentType, *a)
}

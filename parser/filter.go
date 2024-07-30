package parser

import (
	"reflect"
)

var (
	AlwaysTrue  Predicate[Event] = PredicateFunc[Event](TrueFn[Event])
	AlwaysFalse Predicate[Event] = PredicateFunc[Event](FalseFn[Event])
)

func TrueFn[T any](T) bool {
	return true
}

func FalseFn[T any](T) bool {
	return false
}

type EventFilter interface {
	GetPredicate(metadata *ClassMetadata) Predicate[Event]
}

type Predicate[T any] interface {
	Test(t T) bool
}

type PredicateFunc[T any] func(t T) bool

func (p PredicateFunc[T]) Test(t T) bool {
	return p(t)
}

func (p PredicateFunc[T]) Equals(other PredicateFunc[T]) bool {
	return reflect.ValueOf(p).Pointer() == reflect.ValueOf(other).Pointer()
}

func IsAlwaysTrue(p Predicate[Event]) bool {
	if pf, ok := p.(PredicateFunc[Event]); ok {
		return pf.Equals(AlwaysTrue.(PredicateFunc[Event]))
	}
	return false
}

func IsAlwaysFalse(p Predicate[Event]) bool {
	if pf, ok := p.(PredicateFunc[Event]); ok {
		return pf.Equals(AlwaysFalse.(PredicateFunc[Event]))
	}
	return false
}

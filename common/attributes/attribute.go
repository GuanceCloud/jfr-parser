package attributes

import (
	"fmt"
	"github.com/grafana/jfr-parser/common/units"
	"github.com/grafana/jfr-parser/parser"
	"log/slog"
	"reflect"
)

var (
	GcWhen    = Attr[string]("when", "When", "java.lang.String", "")
	Blocking  = Attr[bool]("blocking", "Blocking", "boolean", "Whether the thread calling the vm operation was blocked or not")
	Safepoint = Attr[bool]("safepoint", "Safepoint", "boolean", "Whether the vm operation occured at a safepoint or not")
)

var (
	StartTime           = Attr[units.IQuantity]("startTime", "Start Time", "long", "")
	EventStacktrace     = Attr[*parser.StackTrace]("stackTrace", "Stack Trace", "jdk.types.StackTrace", "")
	EventThread         = Attr[*parser.Thread]("eventThread", "Thread", "java.lang.Thread", "The thread in which the event occurred")
	ThreadStat          = Attr[string]("state", "Thread State", "jdk.types.ThreadState", "")
	CpuSamplingInterval = Attr[units.IQuantity]("cpuInterval", "CPU Sampling Interval", "", "")
	SettingName         = Attr[string]("name", "Setting Name", "java.lang.String", "")
	SettingValue        = Attr[string]("value", "Setting Value", "java.lang.String", "")
	SettingUnit         = Attr[string]("unit", "Setting Unit", "java.lang.String", "")
	DatadogEndpoint     = Attr[string]("endpoint", "Endpoint", "java.lang.String", "")
)

type Attribute[T any] struct {
	Name        string // unique identifier for attribute
	Label       string // human-readable name
	ClassName   string
	Description string
}

func Attr[T any](name, label, className, description string) *Attribute[T] {
	return &Attribute[T]{
		Name:        name,
		Label:       label,
		ClassName:   className,
		Description: description,
	}
}

func AttrSimple[T any](name, className string) *Attribute[T] {
	return &Attribute[T]{
		Name:      name,
		ClassName: className,
	}
}

func AttrNoDesc[T any](name, label, className string) *Attribute[T] {
	return &Attribute[T]{
		Name:      name,
		Label:     label,
		ClassName: className,
	}
}

func (a *Attribute[T]) GetValue(event *parser.GenericEvent) (T, error) {
	var t T
	attr, ok := event.Attributes[a.Name]
	if !ok {
		return t, fmt.Errorf("attribute name [%s] is not found in the event", a.Name)
	}

	if x, ok := attr.(T); ok {
		return x, nil
	}

	attrValue := reflect.ValueOf(attr)
	attrType := attrValue.Type()
	tValue := reflect.ValueOf(&t).Elem()
	tType := tValue.Type()

	if attrType.ConvertibleTo(tType) {
		// t = t(attr)
		tValue.Set(attrValue.Convert(tType))
		return t, nil
	} else if attrValue.Kind() == reflect.Pointer && attrValue.Elem().Type().ConvertibleTo(tType) {
		// t = t(*attr)
		tValue.Set(attrValue.Elem().Convert(tType))
		return t, nil
	} else if tType.Kind() == reflect.Pointer && attrType.ConvertibleTo(tType.Elem()) {
		// t = t(&attr)
		ap := reflect.New(attrType)
		ap.Elem().Set(attrValue)
		if ap.Type().ConvertibleTo(tType) {
			tValue.Set(ap.Convert(tType))
			return t, nil
		}
	}

	fieldMeta := event.ClassMetadata.GetField(a.Name)
	fieldUnit := fieldMeta.Unit(event.ClassMetadata.ClassMap)

	slog.Debug("unit: ", slog.Bool("is nil", fieldUnit == nil))
	if fieldUnit != nil {
		slog.Debug("unit: ", slog.String("name", fieldUnit.Name))
	}

	if fieldUnit != nil || fieldMeta.TickTimestamp(event.ClassMetadata.ClassMap) {
		var (
			num      units.Numeric
			quantity units.IQuantity
		)

		switch attr.(type) {
		case *parser.Byte, *parser.Short, *parser.Int, *parser.Long:
			if fieldMeta.Unsigned(event.ClassMetadata.ClassMap) {
				var x any
				switch ax := attr.(type) {
				case *parser.Byte:
					x = uint8(*ax)
				case *parser.Short:
					x = uint16(*ax)
				case *parser.Int:
					x = uint32(*ax)
					//case *parser.Long: // parser.Long doesn't support unsigned yet
					//	x = uint64(*ax)
				}
				num = units.I64(reflect.ValueOf(x).Uint())
			} else {
				num = units.I64(reflect.ValueOf(attr).Elem().Int())
			}
		case *parser.Float, *parser.Double:
			num = units.F64(reflect.ValueOf(attr).Elem().Float())
		}

		if fieldMeta.TickTimestamp(event.ClassMetadata.ClassMap) {
			ts := fieldMeta.ChunkHeader.StartTimeNanos + ((num.Int64() - fieldMeta.ChunkHeader.StartTicks) * 1e9 / fieldMeta.ChunkHeader.TicksPerSecond)
			quantity = units.NewIntQuantity(ts, units.UnixNano)
		} else {
			if num.Float() {
				quantity = units.NewFloatQuantity(num.Float64(), fieldUnit)
			} else {
				quantity = units.NewIntQuantity(num.Int64(), fieldUnit)
			}
		}

		if q, ok := quantity.(T); ok {
			return q, nil
		}
	}

	switch any(t).(type) {
	case string:
		s, err := parser.ToString(attr)
		if err != nil {
			return t, fmt.Errorf("unable to resolve string: %w", err)
		}
		reflect.ValueOf(&t).Elem().SetString(s)
		return t, nil
	}

	return t, fmt.Errorf("attribute is not type of %T", t)
}

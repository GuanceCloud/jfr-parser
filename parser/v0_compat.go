package parser

import (
	types2 "github.com/grafana/jfr-parser/common/types"
	v0 "github.com/grafana/jfr-parser/parser/v0"
)

type v0FieldSpec struct {
	Name      string
	ClassName string
}

type v0ClassRegistry struct {
	nextID   int64
	classes  map[string]*ClassMetadata
	classMap ClassMap
}

func newV0ClassRegistry() *v0ClassRegistry {
	return &v0ClassRegistry{
		nextID:   -1,
		classes:  make(map[string]*ClassMetadata),
		classMap: make(ClassMap),
	}
}

func (r *v0ClassRegistry) ensureClass(name, superType string) *ClassMetadata {
	if classMeta, ok := r.classes[name]; ok {
		if classMeta.SuperType == "" && superType != "" {
			classMeta.SuperType = superType
		}
		return classMeta
	}

	classMeta := &ClassMetadata{
		ID:        r.nextID,
		Name:      name,
		SuperType: superType,
	}
	r.nextID--
	r.classes[name] = classMeta
	r.classMap[classMeta.ID] = classMeta
	return classMeta
}

func (r *v0ClassRegistry) finalize() ClassMap {
	for _, classMeta := range r.classMap {
		classMeta.ClassMap = r.classMap
	}
	return r.classMap
}

func buildV0EventClassMetadata(registry *v0ClassRegistry, eventType string, fields []v0FieldSpec) *ClassMetadata {
	classMeta := registry.ensureClass(eventType, EventSuperType)
	classMeta.Fields = make([]*FieldMetadata, 0, len(fields))

	for _, field := range fields {
		fieldClass := registry.ensureClass(field.ClassName, "")
		classMeta.Fields = append(classMeta.Fields, &FieldMetadata{
			ClassID: fieldClass.ID,
			Name:    field.Name,
		})
	}

	return classMeta
}

func buildV0FieldSpecs(etd *v0.EventTypeDescriptor, contentTypes map[int]*v0.ContentTypeDescriptor) ([]v0FieldSpec, error) {
	fields := make([]v0FieldSpec, 0, len(etd.ValueDescriptors)+4)
	fields = append(fields, v0FieldSpec{Name: "startTime", ClassName: string(types2.Long)})
	if etd.HasStartTime {
		fields = append(fields, v0FieldSpec{Name: "duration", ClassName: string(types2.Long)})
	}
	if etd.HasThread {
		fields = append(fields, v0FieldSpec{Name: "eventThread", ClassName: string(types2.Thread)})
	}
	if etd.CanHaveStacktrace {
		fields = append(fields, v0FieldSpec{Name: "stackTrace", ClassName: string(types2.StackTrace)})
	}

	payloadFields, err := flattenV0Fields(etd.ValueDescriptors, contentTypes, "")
	if err != nil {
		return nil, err
	}
	fields = append(fields, payloadFields...)
	return fields, nil
}

func flattenV0Fields(vds []*v0.ValueDescriptor, contentTypes map[int]*v0.ContentTypeDescriptor, prefix string) ([]v0FieldSpec, error) {
	fields := make([]v0FieldSpec, 0, len(vds))
	for _, descriptor := range vds {
		name := descriptor.Identifier
		if prefix != "" {
			name = prefix + "." + descriptor.Identifier
		}

		if descriptor.DataType == v0.Struct {
			children, err := descriptor.Children()
			if err != nil {
				return nil, err
			}
			childFields, err := flattenV0Fields(children, contentTypes, name)
			if err != nil {
				return nil, err
			}
			fields = append(fields, childFields...)
			continue
		}

		fields = append(fields, v0FieldSpec{
			Name:      name,
			ClassName: v0FieldClass(descriptor, contentTypes),
		})
	}
	return fields, nil
}

func v0FieldClass(vd *v0.ValueDescriptor, contentTypes map[int]*v0.ContentTypeDescriptor) string {
	if contentType, ok := contentTypes[vd.ContentType]; ok {
		if className := v0ContentTypeClassName(contentType.Name); className != "" {
			return className
		}
	}

	switch vd.DataType {
	case v0.Boolean:
		return string(types2.Boolean)
	case v0.Byte, v0.U1:
		return string(types2.Byte)
	case v0.Short, v0.U2:
		return string(types2.Short)
	case v0.Integer, v0.U4:
		return string(types2.Int)
	case v0.Long, v0.U8:
		return string(types2.Long)
	case v0.Float:
		return string(types2.Float)
	case v0.Double:
		return string(types2.Double)
	case v0.String, v0.UTF8:
		return string(types2.String)
	case v0.Array, v0.StructArray, v0.Struct:
		return "java.lang.Object"
	default:
		return "java.lang.Object"
	}
}

func v0ContentTypeClassName(name string) string {
	switch name {
	case "Thread", "JavaThread", "VMThread":
		return string(types2.Thread)
	case "StackTrace":
		return string(types2.StackTrace)
	case "Class":
		return string(types2.Class)
	case "Method":
		return string(types2.Method)
	case "ThreadGroup":
		return string(types2.ThreadGroup)
	case "ThreadState":
		return string(types2.ThreadState)
	case "GCName":
		return string(types2.String)
	case "GCCause":
		return string(types2.String)
	case "GCWhen":
		return string(types2.String)
	case "G1YCType":
		return string(types2.G1YCType)
	case "GCThresholdUpdater":
		return string(types2.GCThresholdUpdater)
	case "ReferenceType":
		return string(types2.ReferenceType)
	case "MetadataType":
		return string(types2.MetadataType)
	case "MetaspaceObjectType":
		return string(types2.MetaspaceObjectType)
	case "NARROW_OOP_MODE":
		return string(types2.NarrowOopMode)
	case "VMOperationType":
		return string(types2.VMOperationType)
	case "CompilerPhaseType":
		return string(types2.CompilerPhaseType)
	case "FlagValueOrigin":
		return string(types2.FlagValueOrigin)
	case "FrameType":
		return string(types2.FrameType)
	case "UTFConstant":
		return string(types2.String)
	default:
		return ""
	}
}

func normalizeV0Value(value any) any {
	switch v := value.(type) {
	case *v0.LabeledIdentifier:
		return v.Name
	case []any:
		values := make([]any, 0, len(v))
		for _, item := range v {
			values = append(values, normalizeV0Value(item))
		}
		return values
	default:
		return value
	}
}

func wrapV0Value(className string, value any) ParseResolvable {
	value = normalizeV0Value(value)
	if value == nil {
		return nil
	}

	switch className {
	case string(types2.Thread):
		if name, ok := value.(string); ok {
			return &Thread{
				OsName:   name,
				JavaName: name,
			}
		}
	case string(types2.StackTrace):
		return &StackTrace{}
	}

	return WrapRawValue(value)
}

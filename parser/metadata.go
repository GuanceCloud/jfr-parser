package parser

import (
	"fmt"
	"strconv"

	"github.com/grafana/jfr-parser/parser/types/def"
)

type Element interface {
	SetAttribute(key, value string) error
	AppendChild(name string) Element
}

type AnnotationMetadata struct {
	Class  int64
	Values map[string]string
}

func (a *AnnotationMetadata) SetAttribute(key, value string) (err error) {
	switch key {
	case "class":
		a.Class, err = strconv.ParseInt(value, 10, 64)
	default:
		if a.Values == nil {
			a.Values = make(map[string]string)
		}
		a.Values[key] = value
	}
	return err
}

func (a *AnnotationMetadata) AppendChild(string) Element { return nil }

// TODO: Proper attribute support for SettingMetadata
type SettingMetadata struct {
	Values map[string]string
}

func (s *SettingMetadata) SetAttribute(key, value string) error {
	if s.Values == nil {
		s.Values = make(map[string]string)
	}
	s.Values[key] = value
	return nil
}

func (s *SettingMetadata) AppendChild(string) Element { return nil }

type FieldMetadata struct {
	Class        int64
	Name         string
	ConstantPool bool
	Dimension    int32
	Annotations  []AnnotationMetadata
}

func (f *FieldMetadata) SetAttribute(key, value string) (err error) {
	switch key {
	case "name":
		f.Name = value
	case "class":
		f.Class, err = strconv.ParseInt(value, 10, 64)
	case "constantPool":
		f.ConstantPool, err = parseBool(value)
	case "dimension":
		var n int64
		n, err = strconv.ParseInt(value, 10, 32)
		f.Dimension = int32(n)
	}
	return nil
}

func (f *FieldMetadata) AppendChild(name string) Element {
	switch name {
	case "annotation":
		f.Annotations = append(f.Annotations, AnnotationMetadata{})
		return &f.Annotations[len(f.Annotations)-1]
	}
	return nil
}

type ClassMetadata struct {
	ID          int64
	Name        string
	SuperType   string
	SimpleType  bool
	Fields      []FieldMetadata
	Settings    []SettingMetadata
	Annotations []AnnotationMetadata
}

func (c *ClassMetadata) SetAttribute(key, value string) (err error) {
	switch key {
	case "id":
		c.ID, err = strconv.ParseInt(value, 10, 64)
	case "name":
		c.Name = value
	case "superType":
		c.SuperType = value
	case "simpleType":
		c.SimpleType, err = parseBool(value)
	}
	return err
}

func (c *ClassMetadata) AppendChild(name string) Element {
	switch name {
	case "field":
		c.Fields = append(c.Fields, FieldMetadata{})
		return &c.Fields[len(c.Fields)-1]
	case "setting":
		c.Settings = append(c.Settings, SettingMetadata{})
		return &c.Settings[len(c.Settings)-1]
	case "annotation":
		c.Annotations = append(c.Annotations, AnnotationMetadata{})
		return &c.Annotations[len(c.Annotations)-1]
	}
	return nil
}

type Metadata struct {
	Classes []*ClassMetadata
}

func (m *Metadata) SetAttribute(string, string) error { return nil }

func (m *Metadata) AppendChild(name string) Element {
	switch name {
	case "class":
		m.Classes = append(m.Classes, &ClassMetadata{})
		return m.Classes[len(m.Classes)-1]
	}
	return nil
}

type Region struct {
	Locale        string
	GMTOffset     string
	TicksToMillis string
}

func (m *Region) SetAttribute(key, value string) error {
	switch key {
	case "locale":
		m.Locale = value
	case "gmtOffset":
		// TODO int?
		m.GMTOffset = value
	case "ticksToMillis":
		// TODO int?
		m.TicksToMillis = value
	}
	return nil
}

func (m *Region) AppendChild(string) Element { return nil }

type Root struct {
	Metadata *Metadata
	Region   Region
}

func (r *Root) SetAttribute(string, string) error { return nil }

func (r *Root) AppendChild(name string) Element {
	switch name {
	case "metadata":
		r.Metadata = &Metadata{}
		return r.Metadata
	case "region":
		r.Region = Region{}
		return &r.Region
	}
	return nil
}

type MetadataEvent struct {
	StartTime int64
	Duration  int64
	ID        int64
	Root      Root
}

func (m *MetadataEvent) Parse(r Reader) (err error) {
	if kind, err := r.VarLong(); err != nil {
		return fmt.Errorf("unable to retrieve event type: %w", err)
	} else if kind != 0 {
		return fmt.Errorf("unexpected metadata event type: %d", kind)
	}

	if m.StartTime, err = r.VarLong(); err != nil {
		return fmt.Errorf("unable to parse metadata event's start time: %w", err)
	}
	if m.Duration, err = r.VarLong(); err != nil {
		return fmt.Errorf("unable to parse metadata event's duration: %w", err)
	}
	if m.ID, err = r.VarLong(); err != nil {
		return fmt.Errorf("unable to parse metadata event's ID: %w", err)
	}
	n, err := r.VarInt()
	if err != nil {
		return fmt.Errorf("unable to parse metadata event's number of strings: %w", err)
	}
	// TODO: assert n is small enough
	strings := make([]string, n)
	for i := 0; i < int(n); i++ {
		if strings[i], err = r.String(nil); err != nil {
			return fmt.Errorf("unable to parse metadata event's string: %w", err)
		}
	}

	name, err := parseName(r, strings)
	if err != nil {
		return err
	}
	if name != "root" {
		return fmt.Errorf("invalid root element name: %s", name)
	}

	m.Root = Root{}
	if err := parseElement(r, strings, &m.Root); err != nil {
		return fmt.Errorf("unable to parse metadata element tree: %w", err)
	}

	classes := make(map[int64]*ClassMetadata)
	for _, clazz := range m.Root.Metadata.Classes {
		classes[clazz.ID] = clazz
	}

	for _, clazz := range m.Root.Metadata.Classes {
		fmt.Println("metadata class name:", clazz.Name)

		for _, anno := range clazz.Annotations {
			fmt.Println("class annotation class: ", classes[anno.Class].Name, "class annotation value: ", anno.Values)
		}

		for _, field := range clazz.Fields {
			fmt.Println("field name:", field.Name, "field ID:", classes[field.Class].Name, "ConstantPool: ", field.ConstantPool, "Dimension:", field.Dimension)
		}
		fmt.Println("--------------------------------------")

	}

	return nil
}

func (p *Parser) readMeta(pos int) error {
	p.TypeMap.IDMap = make(map[def.TypeID]*def.Class, 43+5)
	p.TypeMap.NameMap = make(map[string]*def.Class, 43+5)

	if err := p.seek(pos); err != nil {
		return err
	}
	sz, err := p.varInt()
	if err != nil {
		return err
	}
	p.metaSize = sz
	_, err = p.varInt()
	if err != nil {
		return err
	}
	_, err = p.varLong()
	if err != nil {
		return err
	}
	_, err = p.varLong()
	if err != nil {
		return err
	}
	_, err = p.varLong()
	if err != nil {
		return err
	}
	nstr, err := p.varInt()
	if err != nil {
		return err
	}
	strings := make([]string, nstr)
	for i := 0; i < int(nstr); i++ {
		strings[i], err = p.string()
		if err != nil {
			return err
		}
	}

	e, err := p.readElement(strings, false)
	if err != nil {
		return err
	}
	if e.name != "root" {
		return fmt.Errorf("expected root element, got %s", e.name)
	}
	for i := 0; i < e.childCount; i++ {
		meta, err := p.readElement(strings, false)
		if err != nil {
			return err
		}
		//fmt.Println(meta.name)
		switch meta.name {
		case "metadata":
			for j := 0; j < meta.childCount; j++ {
				classElement, err := p.readElement(strings, true)

				if err != nil {
					return err
				}
				cls, err := def.NewClass(classElement.attr, classElement.childCount)
				if err != nil {
					return err
				}

				for k := 0; k < classElement.childCount; k++ {
					field, err := p.readElement(strings, true)
					if err != nil {
						return err
					}
					if field.name == "field" {
						f, err := def.NewField(field.attr)
						if err != nil {
							return err
						}
						cls.Fields = append(cls.Fields, f)
					}
					for l := 0; l < field.childCount; l++ {
						_, err := p.readElement(strings, false)
						if err != nil {
							return err
						}
					}

				}
				//fmt.Println(cls.String())
				p.TypeMap.IDMap[cls.ID] = cls
				p.TypeMap.NameMap[cls.Name] = cls

			}
		case "region":
			break
		default:
			return fmt.Errorf("unexpected element %s", meta.name)
		}
	}
	if err := p.checkTypes(); err != nil {
		return err
	}
	return nil
}
func parseElement(r Reader, s []string, e Element) error {
	n, err := r.VarInt()
	if err != nil {
		return fmt.Errorf("unable to parse attribute count: %w", err)
	}
	// TODO: assert n is small enough
	for i := 0; i < int(n); i++ {
		k, err := parseName(r, s)
		if err != nil {
			return fmt.Errorf("unable to parse attribute key: %w", err)
		}
		v, err := parseName(r, s)
		if err != nil {
			return fmt.Errorf("unable to parse attribute value: %w", err)
		}
		if err := e.SetAttribute(k, v); err != nil {
			return fmt.Errorf("unable to set element attribute: %w", err)
		}
	}
	n, err = r.VarInt()
	if err != nil {
		return fmt.Errorf("unable to parse element count: %w", err)
	}
	// TODO: assert n is small enough
	for i := 0; i < int(n); i++ {
		name, err := parseName(r, s)
		if err != nil {
			return fmt.Errorf("unable to parse element name: %w", err)
		}
		child := e.AppendChild(name)
		if child == nil {
			return fmt.Errorf("unexpected child in metadata event: %s", name)
		}
		parseElement(r, s, child)
	}
	return nil
}

func (p *Parser) readElement(strings []string, needAttributes bool) (element, error) {
	iname, err := p.varInt()
	if err != nil {
		return element{}, err
	}
	if iname < 0 || int(iname) >= len(strings) {
		return element{}, def.ErrIntOverflow
	}
	name := strings[iname]
	attributeCount, err := p.varInt()
	if err != nil {
		return element{}, err
	}
	var attributes map[string]string
	if needAttributes {
		attributes = make(map[string]string, attributeCount)
	}
	for i := 0; i < int(attributeCount); i++ {
		attributeName, err := p.varInt()
		if err != nil {
			return element{}, err
		}
		if attributeName < 0 || int(attributeName) >= len(strings) {
			return element{}, def.ErrIntOverflow
		}
		attributeValue, err := p.varInt()
		if err != nil {
			return element{}, err
		}
		if attributeValue < 0 || int(attributeValue) >= len(strings) {
			return element{}, def.ErrIntOverflow
		}
		if needAttributes {
			attributes[strings[attributeName]] = strings[attributeValue]
		} else {
			//fmt.Printf("                              >>> skipping attribute %s=%s\n", strings[attributeName], strings[attributeValue])
		}
	}

	childCount, err := p.varInt()
	if err != nil {
		return element{}, err
	}
	return element{
		name:       name,
		attr:       attributes,
		childCount: int(childCount),
	}, nil

}

func parseName(r Reader, s []string) (string, error) {
	n, err := r.VarInt()
	if err != nil {
		return "", fmt.Errorf("unable to parse string name index: %w", err)
	}
	if int(n) >= len(s) {
		return "", fmt.Errorf("invalid name index %d, only %d names available", n, len(s))
	}
	return s[int(n)], nil
}

func parseBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	return false, fmt.Errorf("unable to parse '%s' as boolean", s)
}

type element struct {
	name       string
	attr       map[string]string
	childCount int
}

package v0

import "fmt"

var (
	Boolean     = &DataType{Name: "Boolean", Size: 1, Numeric: false, Textual: false}
	Byte        = &DataType{Name: "Byte", Size: 1, Numeric: true, Textual: false}
	U1          = &DataType{Name: "U1", Size: 1, Numeric: true, Textual: false} // uint8
	Short       = &DataType{Name: "Short", Size: 2, Numeric: true, Textual: false}
	U2          = &DataType{Name: "U2", Size: 2, Numeric: true, Textual: false} // uint16
	Integer     = &DataType{Name: "Integer", Size: 4, Numeric: true, Textual: false}
	U4          = &DataType{Name: "U4", Size: 4, Numeric: true, Textual: false} // uint32
	Long        = &DataType{Name: "Long", Size: 8, Numeric: true, Textual: false}
	U8          = &DataType{Name: "U8", Size: 8, Numeric: true, Textual: false} // uint64
	Float       = &DataType{Name: "Float", Size: 4, Numeric: true, Textual: false}
	Double      = &DataType{Name: "Double", Size: 8, Numeric: true, Textual: false}
	UTF8        = &DataType{Name: "UTF8", Size: 0, Numeric: false, Textual: true}
	String      = &DataType{Name: "String", Size: 0, Numeric: false, Textual: true}
	Array       = &DataType{Name: "Array", Size: 0, Numeric: false, Textual: false}
	Struct      = &DataType{Name: "Struct", Size: 0, Numeric: false, Textual: false}
	StructArray = &DataType{Name: "StructArray", Size: 0, Numeric: false, Textual: false}
)

var DataTypeEnum = []*DataType{
	Boolean,
	Byte,
	U1,
	Short,
	U2,
	Integer,
	U4,
	Long,
	U8,
	Float,
	Double,
	UTF8,
	String,
	Array,
	Struct,
	StructArray,
}

type ProducerParser struct {
}

type ProducerDescriptor struct {
	/**
	private final String name;
	private final String description;

	private final String uriString;
	private final int identifier;
	private final EventTypeDescriptor[] eventTypeDescriptors;
	private final ContentTypeDescriptor[] contentTypes;
	*/
	Name        string
	Description string
	UriString   string
	Identifier  int

	EventTypeDescriptors   []*EventTypeDescriptor
	ContentTypeDescriptors []*ContentTypeDescriptor
}

type EventTypeDescriptor struct {
	/*
		final private int identifier;
		final private String label;
		final private boolean hasStartTime;
		final private boolean hasThread;
		final private boolean canHaveStacktrace;
		final private boolean isRequstable;
		final private ValueDescriptor[] dataStructure;
		final private String description;
		final private String path;
	*/

	Identifier        int
	Label             string
	HasStartTime      bool
	HasThread         bool
	CanHaveStacktrace bool
	IsRequestAble     bool
	ValueDescriptors  []*ValueDescriptor
	Description       string
	Path              string
}

type ValueDescriptor struct {
	/**
	private final String identifier;
	private final String name;
	private final String description;
	private final Transition transition;
	private final DataType dataType;
	private final String relationalKey;
	private final int contentType;
	private final DataStructure[] dataStructures;
	private final int structureIndex;
	*/
	Identifier     string
	Name           string
	Description    string
	Transition     int
	DataType       *DataType
	RelationalKey  string
	ContentType    int
	DataStructures *DataStructureList
	StructureIndex int
}

func (v *ValueDescriptor) ArrayElemDatatype() (*DataType, error) {
	if v.DataType != Array {
		return nil, fmt.Errorf("only Array type can get element datatype, current datatype: [%s]", v.DataType.Name)
	}
	return DataTypeEnum[v.StructureIndex], nil
}

func (v *ValueDescriptor) Children() ([]*ValueDescriptor, error) {
	if v.DataType != Struct && v.DataType != StructArray {
		return nil, fmt.Errorf("this datatype [%s] doesn't contain children dataStructure", v.DataType.Name)
	}
	return v.DataStructures.Get(int64(v.StructureIndex)).ValueDescriptors, nil
}

func InConstantPool(dataType *DataType, contentType int) bool {
	if !dataType.Numeric {
		return false
	}
	if notInConstantPoolTypes[contentType] {
		return false
	}
	return true
}

func (v *ValueDescriptor) InConstantPool() bool {
	return InConstantPool(v.DataType, v.ContentType)
}

func (v *ValueDescriptor) Desc(indent string) {
	if v.DataType != Array {
		fmt.Printf("%sidentifier: %s | Name: %s | dataType: %s | contentType: %d | structureIndex: %d | is in constant pool: %t\n",
			indent, v.Identifier, v.Name, v.DataType.Name, v.ContentType, v.StructureIndex, v.InConstantPool())
	} else {
		elemType, err := v.ArrayElemDatatype()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%sidentifier: %s | Name: %s | dataType: %s | contentType: %d | structureIndex: %d | array element type: %s | array element is in constant pool: %t\n",
			indent, v.Identifier, v.Name, v.DataType.Name, v.ContentType, v.StructureIndex, elemType.Name, InConstantPool(elemType, v.ContentType))
	}

	switch v.DataType {
	case Struct, StructArray:
		children, err := v.Children()
		if err != nil {
			panic(err)
		}

		for _, child := range children {
			child.Desc(indent + "    ")
		}
	}
}

type DataStructureList []*DataStructure

func NewDataStructureList(capacity int64) *DataStructureList {
	l := make(DataStructureList, 0, capacity)
	return &l
}

func (d *DataStructureList) Len() int {
	return len(*d)
}

func (d *DataStructureList) Add(ds *DataStructure) {
	*d = append(*d, ds)
}

func (d *DataStructureList) Get(i int64) *DataStructure {
	return (*d)[i]
}

type DataStructure struct {
	ValueDescriptors []*ValueDescriptor
}

type DataType struct {
	Name    string
	Size    int
	Numeric bool
	Textual bool
}

var primitiveDataTypes map[string]bool

func (d *DataType) IsPrimitive() bool {
	if primitiveDataTypes == nil {
		primitiveDataTypes = make(map[string]bool)
		for _, dt := range DataTypeEnum {
			if dt == Array {
				break
			}
			primitiveDataTypes[dt.Name] = true
		}
	}
	return primitiveDataTypes[d.Name]
}

type ContentTypeDescriptor struct {
	/**
	private final int contentTypeId;
	private final String name;
	private final String description;
	private final DataType dataType;
	private final ValueDescriptor[] dataStructure;
	*/
	ContentTypeId    int
	Name             string
	Description      string
	DataType         *DataType
	ValueDescriptors []*ValueDescriptor
}

func (c *ContentTypeDescriptor) InConstantPool() bool {
	return InConstantPool(c.DataType, c.ContentTypeId)
}

func (c *ContentTypeDescriptor) Desc() {
	fmt.Printf("\ncontent type id: %d | content type Name: %s | content type DataType: %s | is in constant pool: %t\n",
		c.ContentTypeId, c.Name, c.DataType.Name, c.InConstantPool())
	for _, vd := range c.ValueDescriptors {
		vd.Desc("  ")
	}
}

package parser

import (
	"strings"
	"testing"
)

func TestClassMetadata_Category(t *testing.T) {
	chunks, err := ParseFile("./testdata/prof.jfr")
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {

		for _, class := range chunk.Metadata.ClassMap {
			t.Logf("class id: %d, class name: %s, simple type: %t, super type: %s, class category: %s",
				class.ID, class.Name, class.SimpleType, class.SuperType, strings.Join(class.Category(), " --> "))

			for _, field := range class.Fields {
				t.Logf("field name: %s, field label: %s, field description: %s, field isArray: %t, field inConstantPool: %t, field class: %s",
					field.Name, field.Label(chunk.Metadata.ClassMap), field.Description(chunk.Metadata.ClassMap),
					field.IsArray(), field.ConstantPool, chunk.Metadata.ClassMap[field.ClassID].Name)
			}
		}
	}

}

package parser

import (
	types2 "github.com/grafana/jfr-parser/common/types"
	"testing"
)

func TestClassMetadata_Category(t *testing.T) {
	chunks, err := ParseFile("./testdata/prof.jfr")
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {
		chunk.ShowClassMeta(types2.DatadogProfilerConfig)
	}

}

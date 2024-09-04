package filters

import (
	"github.com/grafana/jfr-parser/common/attributes"
	"github.com/grafana/jfr-parser/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAlwaysTrue(t *testing.T) {
	var a, b parser.PredicateFunc
	assert.True(t, a.Equals(b))

	var c = parser.AlwaysTrue.(parser.PredicateFunc)
	var d = parser.PredicateFunc(parser.TrueFn)

	assert.True(t, c.Equals(d))

	assert.True(t, parser.IsAlwaysTrue(parser.AlwaysTrue))
	assert.True(t, parser.IsAlwaysFalse(parser.AlwaysFalse))
	assert.False(t, parser.IsAlwaysTrue(parser.AlwaysFalse))
	assert.False(t, parser.IsAlwaysFalse(parser.AlwaysTrue))
	assert.False(t, parser.IsAlwaysFalse(a))

}

func TestAndAlways(t *testing.T) {
	var a parser.Predicate[parser.Event] = parser.PredicateFunc(func(ge parser.Event) bool {
		_, ok := ge.(*parser.GenericEvent)
		return ok
	})

	var isNil parser.Predicate[parser.Event] = parser.PredicateFunc(func(e parser.Event) bool {
		return e == nil
	})

	ap1 := AndAlways(parser.AlwaysTrue, parser.AlwaysTrue, parser.AlwaysTrue)
	assert.True(t, parser.IsAlwaysTrue(ap1))

	ap2 := AndAlways(a, isNil, parser.AlwaysFalse)
	assert.True(t, parser.IsAlwaysFalse(ap2))

	ap3 := NotAlways(parser.AlwaysTrue)
	assert.True(t, parser.IsAlwaysFalse(ap3))

	ap4 := NotAlways(parser.AlwaysFalse)
	assert.True(t, parser.IsAlwaysTrue(ap4))

	ap5 := AndAlways(parser.AlwaysTrue, parser.AlwaysFalse, parser.AlwaysTrue)
	assert.True(t, parser.IsAlwaysFalse(ap5))

	ap6 := OrAlways(parser.AlwaysFalse, parser.AlwaysFalse, parser.AlwaysFalse)
	assert.True(t, parser.IsAlwaysFalse(ap6))

	ap7 := OrAlways(parser.AlwaysFalse, parser.AlwaysTrue, parser.AlwaysFalse)
	assert.True(t, parser.IsAlwaysTrue(ap7))

	ap8 := AndAlways(isNil, parser.AlwaysTrue, parser.AlwaysTrue, parser.AlwaysTrue)
	assert.False(t, parser.IsAlwaysTrue(ap8))
	assert.False(t, parser.IsAlwaysFalse(ap8))
	assert.True(t, ap8.Test(nil))

	ap9 := NotAlways(isNil)
	assert.False(t, parser.IsAlwaysTrue(ap9))
	assert.False(t, parser.IsAlwaysFalse(ap9))
	assert.False(t, ap9.Test(nil))

}

func TestAttributes(t *testing.T) {
	chunks, err := parser.ParseFile("testdata/prof.jfr")
	if err != nil {
		t.Fatal(err)
	}
	for _, chunk := range chunks {
		//for _, collection := range chunk.Apply(VmInfo) {
		//	for _, event := range collection.Events {
		//		startTime, err := attributes.JVMStartTime.GetValue(event)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		t.Logf("iquantity: %d, unit: %v", startTime.IntValue(), *startTime.Unit())
		//		st, err := units.ToTime(startTime)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		t.Log("jvm start at: ", st)
		//	}
		//}

		//for _, collection := range chunk.Apply(DatadogExecutionSample) {
		//	for _, event := range collection.Events {
		//		weight, err := attributes.SampleWeight.GetValue(event)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		t.Logf("weight: %d", weight)
		//	}
		//}

		for _, event := range chunk.Apply(DatadogProfilerConfig) {
			cpu, err := attributes.CpuSamplingInterval.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}
			wall, err := attributes.WallSampleInterval.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("cpu: ", cpu.String())
			t.Log("wall: ", wall.String())
		}
	}
}

func TestTypes(t *testing.T) {
	chunks, err := parser.ParseFile("./testdata/prof.jfr")
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {

		for _, event := range chunk.ChunkEvents.Apply(DatadogExecutionSample) {

			stackTrace, err := attributes.EventStacktrace.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			for _, frame := range stackTrace.Frames {

				if frame.Type != nil {
					t.Logf("frameType description: %s", frame.Type.Description)
				}

				if frame.Method != nil {
					t.Logf("method name: %s, method class: %s", frame.Method.Name.String, frame.Method.Type.Name.String)
				}

			}
		}
	}

	attr := attributes.AttrSimple[*parser.InflateCause]("cause", "jdk.types.InflateCause")

	for _, chunk := range chunks {
		for _, event := range chunk.ChunkEvents.Apply(JavaMonitorInflate) {

			av, err := attr.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("av.cause: %s", av.Cause)
		}
	}

	for _, chunk := range chunks {
		for _, event := range chunk.ChunkEvents.Apply(DatadogEndpoint) {
			ep, err := attributes.DatadogEndpoint.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("endpoint: %s", ep)
		}
	}
}

func TestParseFile(t *testing.T) {
	cks, err := parser.ParseFile("testdata/corrupt.jfr")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(cks))
}

func TestParseLZ4(t *testing.T) {
	chunks, err := parser.ParseFile("./testdata/ddtrace.jfr.lz4")

	if err != nil {
		t.Fatalf("Unable to parse jfr file: %s", err)
	}

	for _, chunk := range chunks {

		for _, event := range chunk.Apply(FilterExecutionSample) {

			for _, field := range event.ClassMetadata.Fields {
				t.Logf("field: %s, class: %s, unit: %v",
					field.Name, event.ClassMetadata.ClassMap[field.ClassID].Name,
					field.Unit(event.ClassMetadata.ClassMap))

				for _, annotation := range field.Annotations {
					t.Logf("class: %s, values: %v", event.ClassMetadata.ClassMap[annotation.ClassID].Name, annotation.Values)
				}
			}

			startTime, err := attributes.StartTime.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("%d %s", startTime.IntValue(), startTime.Unit().Name)

		}

		for _, event := range chunk.Apply(DatadogProfilerSetting) {

			name, err := attributes.SettingName.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			value, err := attributes.SettingValue.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			unit, err := attributes.SettingUnit.GetValue(event)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("setting name: %q, setting value: %q, setting unit: %q", name, value, unit)

			//for key, attr := range ge.Attributes {
			//	t.Logf("%s: %+#v (%T)\n", key, attr, attr)
			//}
			t.Logf("\n")
		}
	}
}

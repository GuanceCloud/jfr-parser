package filters

import (
	"github.com/grafana/jfr-parser/common/attributes"
	"github.com/grafana/jfr-parser/parser"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strings"
	"testing"
)

type User struct {
	Name  string
	Age   int
	Score int
}

func isAdult(user *User) bool {
	return user.Age >= 18
}

func hasLastName(user *User) bool {
	return len(strings.Split(user.Name, " ")) >= 2
}

func excellent(u *User) bool {
	return u.Score >= 95
}

func TestAnd(t *testing.T) {
	a := parser.PredicateFunc[*User](isAdult)
	l := parser.PredicateFunc[*User](hasLastName)
	e := parser.PredicateFunc[*User](excellent)

	ap := And[*User](a, l, e)

	assert.True(t, ap.Test(&User{Name: "Mary William", Age: 19, Score: 99}))
	assert.False(t, ap.Test(&User{Name: "Tom Crux", Age: 19}))

	op := Or[*User](a, l, e)

	assert.True(t, op.Test(&User{
		Age:   17,
		Name:  "Jerry",
		Score: 99,
	}))

	assert.False(t, op.Test(&User{
		Age:   17,
		Name:  "Jerry",
		Score: 60,
	}))
}

func TestIsAlwaysTrue(t *testing.T) {
	var a, b parser.PredicateFunc[parser.Event]
	assert.True(t, a.Equals(b))

	var c = parser.AlwaysTrue.(parser.PredicateFunc[parser.Event])
	var d = parser.PredicateFunc[parser.Event](parser.TrueFn[parser.Event])

	assert.True(t, c.Equals(d))

	assert.True(t, parser.IsAlwaysTrue(parser.AlwaysTrue))
	assert.True(t, parser.IsAlwaysFalse(parser.AlwaysFalse))
	assert.False(t, parser.IsAlwaysTrue(parser.AlwaysFalse))
	assert.False(t, parser.IsAlwaysFalse(parser.AlwaysTrue))
	assert.False(t, parser.IsAlwaysFalse(a))

}

func TestAndAlways(t *testing.T) {
	var a parser.Predicate[parser.Event] = parser.PredicateFunc[parser.Event](func(ge parser.Event) bool {
		_, ok := ge.(*parser.GenericEvent)
		return ok
	})

	var isNil parser.Predicate[parser.Event] = parser.PredicateFunc[parser.Event](func(e parser.Event) bool {
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

func TestTypes(t *testing.T) {
	chunks, err := parser.ParseFile("./testdata/prof.jfr")
	if err != nil {
		t.Fatal(err)
	}

	for _, chunk := range chunks {

		for _, collection := range chunk.Events.Apply(DatadogExecutionSample) {

			for _, event := range collection.EventList {

				stackTrace, err := attributes.EventStacktrace.GetValue(event.(*parser.GenericEvent))
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
	}

	attr := attributes.AttrSimple[*parser.InflateCause]("cause", "jdk.types.InflateCause")

	for _, chunk := range chunks {
		for className, collection := range chunk.Events.Apply(JavaMonitorInflate) {
			t.Log("class: ", className)

			for _, event := range collection.EventList {
				av, err := attr.GetValue(event.(*parser.GenericEvent))
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("av.cause: %s", av.Cause)
			}
		}
	}

	for _, chunk := range chunks {
		for className, collection := range chunk.Events.Apply(DatadogEndpoint) {
			t.Logf("class name: %s", className)

			for _, event := range collection.EventList {
				ep, err := attributes.DatadogEndpoint.GetValue(event.(*parser.GenericEvent))
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("endpoint: %s", ep)
			}
		}
	}
}

func TestParseLZ4(t *testing.T) {
	chunks, err := parser.ParseFile("./testdata/ddtrace.jfr.lz4")

	if err != nil {
		t.Fatalf("Unable to parse jfr file: %s", err)
	}

	for _, chunk := range chunks {

		for _, ec := range chunk.Events.Apply(FilterExecutionSample) {

			for _, field := range ec.ClassMetadata.Fields {
				t.Logf("field: %s, class: %s, unit: %v",
					field.Name, ec.ClassMetadata.ClassMap[field.ClassID].Name,
					field.Unit(ec.ClassMetadata.ClassMap))

				for _, annotation := range field.Annotations {
					t.Logf("class: %s, values: %v", ec.ClassMetadata.ClassMap[annotation.ClassID].Name, annotation.Values)
				}
			}

			t.Log(len(ec.EventList))

			for _, event := range ec.EventList {
				startTime, err := attributes.StartTime.GetValue(event.(*parser.GenericEvent))
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%d %s", startTime.IntValue(), startTime.Unit().Name)
			}
		}

		for _, ec := range chunk.Events.Apply(DatadogProfilerSetting) {
			meta := ec.ClassMetadata
			for _, field := range meta.Fields {
				slog.Debug("",
					"field.name", field.Name,
					"field.classID", field.ClassID,
					"filed.className", meta.ClassMap[field.ClassID].Name,
					"field.IsArray", field.IsArray(),
					"field.unit", field.Unit(meta.ClassMap),
					"field.unsigned", field.Unsigned(meta.ClassMap),
					"field label", field.Label(meta.ClassMap),
					"field description", field.Description(meta.ClassMap),
				)
			}

			slog.Warn("", "class name", meta.Name, "events count", slog.IntValue(len(ec.EventList)))

			for _, event := range ec.EventList {
				ge := event.(*parser.GenericEvent)

				name, err := attributes.SettingName.GetValue(ge)
				if err != nil {
					t.Fatal(err)
				}

				value, err := attributes.SettingValue.GetValue(ge)
				if err != nil {
					t.Fatal(err)
				}

				unit, err := attributes.SettingUnit.GetValue(ge)
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
}

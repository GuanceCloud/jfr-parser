package units

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntQuantity_In(t *testing.T) {
	q1 := Nanosecond.IntQuantity(35600_000)
	t.Logf("number: %d, unit: %+#v", q1.IntValue(), *q1.Unit())

	_, err := q1.In(UnixNano)
	assert.NotNil(t, err)
	t.Log(err)

	q2, err := q1.In(Microsecond)
	if err != nil {
		t.Fatal(err)
	}
	iq, ok := q2.(*IntQuantity)
	assert.True(t, ok)
	t.Logf("number: %d, unit: %+#v", iq.IntValue(), *iq.unit)

	q3, err := q1.In(Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	fq, ok := q3.(*FloatQuantity)
	assert.True(t, ok)
	t.Logf("number: %f, unit: %+#v", fq.FloatValue(), *fq.unit)

	q4, err := q1.In(Second)
	if err != nil {
		t.Fatal(err)
	}
	fq, ok = q4.(*FloatQuantity)
	assert.True(t, ok)
	t.Logf("number: %f, unit: %+#v", fq.FloatValue(), *fq.unit)
}

func TestToTime(t *testing.T) {

	timeNow := time.Now()

	q1 := UnixMilli.IntQuantity(timeNow.UnixMilli())
	t.Logf("number: %d, unit: %+#v", q1.IntValue(), *q1.Unit())

	iq, err := q1.In(UnixMicro)
	if err != nil {
		t.Fatal(err)
	}

	tm, err := ToTime(iq)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tm)
	assert.Equal(t, timeNow.Unix(), tm.Unix())
}

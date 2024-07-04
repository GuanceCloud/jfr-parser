package parser

import (
	"testing"
)

func TestIsNilValue(t *testing.T) {

	var (
		err   error
		s     []int
		m     map[string]int
		arr   [2]interface{}
		x     interface{}
		iface interface{} = (*int)(nil)
		i     interface{}
	)

	i = &i

	testcases := []struct {
		name  string
		value interface{}
		isNil bool
	}{
		{
			name:  "nil error",
			value: err,
			isNil: true,
		},
		{
			name:  "nil slice",
			value: s,
			isNil: true,
		},
		{
			name:  "address of nil slice",
			value: &s,
			isNil: true,
		},
		{
			name:  "nil map",
			value: m,
			isNil: true,
		},
		{
			name:  "address of nil map",
			value: &m,
			isNil: true,
		},
		{
			name:  "empty map",
			value: map[string]interface{}{},
			isNil: false,
		},
		{
			name:  "array",
			value: arr,
			isNil: false,
		},
		{
			name:  "nil interface{}",
			value: x,
			isNil: true,
		},
		{
			name:  "nil typed interface{}",
			value: iface,
			isNil: true,
		},
		{
			name:  "address of nil interface{}",
			value: i,
			isNil: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if ok := isNilValue(tc.value); ok != tc.isNil {
				t.Fatalf("%v expected, got %v", tc.isNil, ok)
			}
		})
	}
}

package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceDataFromCtx(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		td   TraceData
	}{
		{
			name: "empty",
			ctx:  context.Background(),
			td:   TraceData{},
		},
		{
			name: "ok",
			ctx:  context.WithValue(context.Background(), traceDataCtxKey{}, TraceData{}),
			td:   TraceData{},
		},
		{
			name: "ok wrong type",
			ctx:  context.WithValue(context.Background(), traceDataCtxKey{}, "string"),
			td:   TraceData{},
		},
		{
			name: "with value",
			ctx: context.WithValue(context.Background(), traceDataCtxKey{}, TraceData{
				"true": TraceValueBool(true),
				"int":  TraceValueInt(2),
				"str":  TraceValueString("mystr"),
			}),
			td: TraceData{
				"true": TraceValueBool(true),
				"int":  TraceValueInt(2),
				"str":  TraceValueString("mystr"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := TraceDataFromCtx(tc.ctx)
			assert.Equal(t, tc.td, actual)
		})
	}
}

type testPayload struct {
	Name string
}

func (t testPayload) Extract() TraceData {
	return TraceData{
		"name": TraceValueString(t.Name),
	}
}

func TestTraceDataAppend(t *testing.T) {
	tests := []struct {
		name      string
		fromTD    TraceData
		extractor TraceDataExtractor
		kv        []TraceKV
		expected  TraceData
	}{
		{
			name: "full with overwrite value",
			fromTD: TraceData{
				"true": TraceValueBool(true),
				"int":  TraceValueInt(2),
				"str":  TraceValueString("mystr"),
			},
			extractor: testPayload{Name: "test name"},
			kv: []TraceKV{
				TraceKVString("str", "someStrValue 2"),
				TraceKVString("someStr", "someStrValue"),
			},
			expected: TraceData{
				"true":    TraceValueBool(true),
				"int":     TraceValueInt(2),
				"str":     TraceValueString("someStrValue 2"),
				"name":    TraceValueString("test name"),
				"someStr": TraceValueString("someStrValue"),
			},
		},
		{
			name:      "all full",
			fromTD:    TraceData{},
			extractor: nil,
			kv:        []TraceKV{},
			expected:  TraceData{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := TraceDataAppend(tc.fromTD, tc.extractor, tc.kv...)
			assert.Equal(t, tc.expected, actual)
		})
	}

}

func TestKV(t *testing.T) {
	boolTrue := true
	strVal := "true"
	intVal := 2
	tests := []struct {
		name     string
		fromKV   TraceKV
		expected TraceKV
	}{
		{
			name:     "bol",
			fromKV:   TraceKVBool("bol", true),
			expected: TraceKV{key: "bol", value: TraceValue{boolValue: &boolTrue}},
		},
		{
			name:     "int",
			fromKV:   TraceKVInt("int", intVal),
			expected: TraceKV{key: "int", value: TraceValue{intValue: &intVal}},
		},
		{
			name:     "str",
			fromKV:   TraceKVString("str", strVal),
			expected: TraceKV{key: "str", value: TraceValue{strValue: &strVal}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.fromKV)
		})
	}

}

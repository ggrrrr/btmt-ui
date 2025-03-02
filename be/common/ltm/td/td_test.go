package td

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttr(t *testing.T) {
	testCases := []struct {
		desc   string
		from   TraceData
		exp    []slog.Attr
		expLen int
	}{
		{
			desc:   "ok empty with attr",
			from:   TraceData{KV: map[string]slog.Value{}},
			exp:    []slog.Attr{},
			expLen: 0,
		},
		{
			desc:   "ok empty",
			from:   TraceData{KV: nil},
			exp:    []slog.Attr{},
			expLen: 0,
		},
		{
			desc: "ok one",
			from: TraceData{map[string]slog.Value{
				"s": slog.StringValue("v1"),
			}},
			exp: []slog.Attr{
				slog.String("s", "v1"),
			},
			expLen: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			res := tc.from.Attr()
			assert.Equal(t, tc.exp, res)
			assert.Equal(t, tc.expLen, tc.from.Len())
		})
	}
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		desc string
		from []slog.Attr
		exp  *TraceData
	}{
		{
			desc: "ok group",
			from: []slog.Attr{slog.Any("g", slog.GroupValue(slog.String("s", "s1")))},
			exp: &TraceData{
				KV: map[string]slog.Value{
					"g": slog.GroupValue(slog.String("s", "s1")),
				},
			},
		},
		{
			desc: "ok",
			from: []slog.Attr{slog.String("str1", "val1"), slog.Int("int1", 123)},
			exp: &TraceData{
				KV: map[string]slog.Value{
					"str1": slog.StringValue("val1"),
					"int1": slog.IntValue(123),
				},
			},
		},
		{
			desc: "empty",
			from: []slog.Attr{},
			exp: &TraceData{
				KV: map[string]slog.Value{},
			},
		},
		{
			desc: "duplicate",
			from: []slog.Attr{slog.String("str1", "val1"), slog.String("str1", "val2")},
			exp: &TraceData{
				KV: map[string]slog.Value{
					"str1": slog.StringValue("val2"),
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			exp := create(tC.from...)
			assert.Equal(t, tC.exp, exp)
		})
	}
}

func TestAppend(t *testing.T) {
	testCases := []struct {
		desc  string
		from  *TraceData
		attrs []slog.Attr
		exp   *TraceData
	}{
		{
			desc:  "from empty",
			from:  &TraceData{},
			attrs: []slog.Attr{slog.String("key1", "val1"), slog.Int("int1", 1)},
			exp: &TraceData{KV: map[string]slog.Value{
				"key1": slog.StringValue("val1"),
				"int1": slog.IntValue(1),
			}},
		},
		{
			desc:  "from not empty",
			from:  &TraceData{KV: map[string]slog.Value{"init": slog.StringValue("first")}},
			attrs: []slog.Attr{slog.String("key1", "val1"), slog.Int("int1", 1)},
			exp: &TraceData{KV: map[string]slog.Value{
				"init": slog.StringValue("first"),
				"key1": slog.StringValue("val1"),
				"int1": slog.IntValue(1),
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.from.append(tC.attrs...)
			assert.Equal(t, tC.exp, tC.from)
		})
	}
}

type myData1 struct {
	s string
}

func (m myData1) Extract() *TraceData {
	return &TraceData{
		KV: map[string]slog.Value{
			"s1": slog.StringValue(m.s),
		},
	}
}

type myData2 struct {
	s string
}

func (m myData2) Extract() *TraceData {
	return &TraceData{
		KV: map[string]slog.Value{
			"s2": slog.StringValue(m.s),
		},
	}
}

type myDataNil struct {
	s string
}

func (m myDataNil) Extract() *TraceData {
	return nil
}

func TestInject(t *testing.T) {
	testCases := []struct {
		desc      string
		from      []TraceDataExtractor
		groupName string
		groupAttr []slog.Attr
		result    *TraceData
	}{
		{
			desc: "ok",
			from: []TraceDataExtractor{&myData1{s: "val1"}},
			result: &TraceData{KV: map[string]slog.Value{
				"s1": slog.StringValue("val1"),
			}},
		},
		{
			desc:   "ok nil",
			from:   []TraceDataExtractor{nil},
			result: &TraceData{KV: map[string]slog.Value{}},
		},
		{
			desc:   "ok extract nil",
			from:   []TraceDataExtractor{myDataNil{}},
			result: &TraceData{KV: map[string]slog.Value{}},
		},
		{
			desc: "ok extract nil 2",
			from: []TraceDataExtractor{myDataNil{}, &myData1{s: "val1"}},
			result: &TraceData{KV: map[string]slog.Value{
				"s1": slog.StringValue("val1"),
			}},
		},
		{
			desc: "ok two injects",
			from: []TraceDataExtractor{&myData1{s: "val1"}, &myData2{s: "val2"}},
			result: &TraceData{KV: map[string]slog.Value{
				"s1": slog.StringValue("val1"),
				"s2": slog.StringValue("val2"),
			}},
		},
		{
			desc:      "ok td and group",
			from:      []TraceDataExtractor{&myData1{s: "val1"}},
			groupName: "g1",
			groupAttr: []slog.Attr{slog.String("b", "s1")},
			result: &TraceData{KV: map[string]slog.Value{
				"s1": slog.StringValue("val1"),
				"g1": slog.GroupValue(slog.String("b", "s1")),
			}},
		},
		{
			desc:      "ok ampty td and group",
			groupName: "g1",
			groupAttr: []slog.Attr{slog.String("b", "ss1")},
			result: &TraceData{KV: map[string]slog.Value{
				"g1": slog.GroupValue(slog.String("b", "ss1")),
			}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := context.Background()
			for _, i := range tc.from {
				ctx = Inject(ctx, i)
			}
			if len(tc.groupAttr) > 0 {
				ctx = InjectGroup(ctx, tc.groupName, tc.groupAttr...)
			}
			td := Extract(ctx)
			if !assert.Equal(t, tc.result, td) {
				fmt.Printf("%v", td)
			}

		})
	}
}

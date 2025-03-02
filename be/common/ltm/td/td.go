package td

import (
	"context"
	"log/slog"
)

type (
	TraceDataExtractor interface {
		Extract() *TraceData
	}

	TraceData struct {
		KV map[string]slog.Value
	}

	tdKeyType struct{}
)

var tdKey tdKeyType = tdKeyType{}

func (td *TraceData) Attr() []slog.Attr {
	if td == nil {
		return []slog.Attr{}
	}
	out := make([]slog.Attr, 0, len(td.KV))

	for k := range td.KV {
		out = append(out, slog.Attr{Key: k, Value: td.KV[k]})
	}
	return out
}

func (td *TraceData) Len() int {
	if td.KV == nil {
		td.KV = make(map[string]slog.Value)
		return 0
	}
	return len(td.KV)
}

func Inject(ctx context.Context, td TraceDataExtractor) context.Context {
	if td == nil {
		return ctx
	}
	c := ctx.Value(tdKey)
	currentTD, ok := c.(*TraceData)
	if !ok {
		return context.WithValue(ctx, tdKey, create(td.Extract().Attr()...))
	}
	currentTD.append(td.Extract().Attr()...)
	return context.WithValue(ctx, tdKey, currentTD)
}

func InjectGroup(ctx context.Context, key string, kv ...slog.Attr) context.Context {
	val := slog.Any(key, slog.GroupValue(kv...))
	c := ctx.Value(tdKey)
	currentTD, ok := c.(*TraceData)
	if !ok {
		return context.WithValue(ctx, tdKey, create(val))
	}
	currentTD.append(val)
	return context.WithValue(ctx, tdKey, currentTD)
}

func Extract(ctx context.Context) *TraceData {
	c := ctx.Value(tdKey)
	currentTD, ok := c.(*TraceData)
	if ok {
		return currentTD
	}
	return &TraceData{KV: make(map[string]slog.Value)}
}

func create(kv ...slog.Attr) *TraceData {
	out := &TraceData{
		KV: map[string]slog.Value{},
	}
	for _, v := range kv {
		out.KV[v.Key] = v.Value
	}
	return out
}

func (td *TraceData) append(kv ...slog.Attr) {
	if td.KV == nil {
		td.KV = make(map[string]slog.Value)
	}
	for i := range kv {
		td.KV[kv[i].Key] = kv[i].Value
	}
}

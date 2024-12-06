package tracedata

import "context"

type (
	TraceData map[string]TraceValue

	TraceDataExtractor interface {
		Extract() TraceData
	}

	traceDataCtxKey struct{}
)

func (td TraceData) Value(key string) TraceValue {
	val, ok := td[key]
	if !ok {
		return TraceValue{}
	}
	return val
}

func TraceDataFromCtx(ctx context.Context) TraceData {
	raw := ctx.Value(traceDataCtxKey{})

	traceData, ok := raw.(TraceData)
	if ok {
		return traceData
	}
	return TraceData{}
}

func TraceDataAppend(data TraceData, payload TraceDataExtractor, kv ...TraceKV) TraceData {
	if payload != nil {
		td := payload.Extract()
		for k, v := range td {
			data[k] = v
		}
	}
	if len(kv) > 0 {
		for _, v := range kv {
			data[v.key] = v.value
		}
	}

	return data
}

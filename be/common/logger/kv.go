package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	TraceData map[string]TraceValue

	TraceDataExtractor interface {
		Extract() TraceData
	}

	TraceValue struct {
		strValue  *string
		boolValue *bool
		intValue  *int
	}

	TraceKV struct {
		key   string
		value TraceValue
	}
)

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

func TraceKVBool(k string, val bool) TraceKV {
	return TraceKV{
		key:   k,
		value: TraceValueBool(val),
	}
}

func TraceKVInt(k string, val int) TraceKV {
	return TraceKV{
		key:   k,
		value: TraceValueInt(val),
	}
}

func TraceKVString(k string, val string) TraceKV {
	return TraceKV{
		key:   k,
		value: TraceValueString(val),
	}
}

func TraceValueBool(val bool) TraceValue {
	return TraceValue{
		boolValue: &val,
	}
}

func TraceValueString(val string) TraceValue {
	return TraceValue{
		strValue: &val,
	}
}

func TraceValueInt(val int) TraceValue {
	return TraceValue{
		intValue: &val,
	}
}

func SpanAttributes(authInfo roles.AuthInfo, td TraceData) []attribute.KeyValue {
	out := []attribute.KeyValue{}
	out = append(out, attribute.String(fmt.Sprintf("%s.auth.subject", spanKeyPrefix), authInfo.Subject))
	out = append(out, attribute.String(fmt.Sprintf("%s.auth.device.info", spanKeyPrefix), authInfo.Device.DeviceInfo))
	out = append(out, attribute.String(fmt.Sprintf("%s.auth.device.addr", spanKeyPrefix), authInfo.Device.RemoteAddr))
	out = append(out, attribute.String(fmt.Sprintf("%s.auth.realm", spanKeyPrefix), string(authInfo.Realm)))

	for k, v := range td {
		key := fmt.Sprintf("%s.data.%s", spanKeyPrefix, k)
		out = append(out, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: v.attributeValue(),
		})
	}
	return out
}

func (v TraceValue) attributeValue() attribute.Value {
	if v.boolValue != nil {
		return attribute.BoolValue(*v.boolValue)
	}
	if v.intValue != nil {
		return attribute.IntValue(*v.intValue)
	}
	if v.strValue != nil {
		return attribute.StringValue(*v.strValue)
	}

	return attribute.StringValue("empty or unknown value")
}

func (v TraceValue) Value() any {
	if v.boolValue != nil {
		return *v.boolValue
	}
	if v.intValue != nil {
		return *v.intValue
	}
	if v.strValue != nil {
		return *v.strValue
	}

	return "empty or unknown value"
}

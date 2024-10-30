package logger

import "go.opentelemetry.io/otel/attribute"

type KV struct {
	key       string
	strValue  *string
	boolValue *bool
	intValue  *int
}

func (kv KV) otelKeyValue() attribute.KeyValue {
	key := kv.key
	if kv.boolValue != nil {
		return AttributeBool(key, *kv.boolValue)
	}
	if kv.strValue != nil {
		return AttributeString(key, *kv.strValue)
	}
	if kv.intValue != nil {
		return AttributeInt(key, *kv.intValue)
	}
	return AttributeString(key, "empty value")

}

func KVBool(key string, val bool) KV {
	return KV{
		key:       key,
		boolValue: &val,
	}
}

func KVString(key string, val string) KV {
	return KV{
		key:      key,
		strValue: &val,
	}
}

func KVInt(key string, val int) KV {
	return KV{
		key:      key,
		intValue: &val,
	}
}

package tracedata

type valueType int

const (
	_ valueType = iota
	valueTypeStr
	valueTypeBool
	valueTypeInt64
	valueTypeUInt64
)

type (
	TraceValue struct {
		t         valueType
		strValue  string
		boolValue bool
		intValue  int64
		uintValue uint64
	}

	TraceKV struct {
		key   string
		value TraceValue
	}
)

func TraceKVBool(k string, val bool) TraceKV {
	return TraceKV{
		key:   k,
		value: ValueBool(val),
	}
}

func TraceKVInt(k string, val int64) TraceKV {
	return TraceKV{
		key:   k,
		value: ValueInt64(val),
	}
}

func TraceKVUInt(k string, val uint64) TraceKV {
	return TraceKV{
		key:   k,
		value: ValueUInt64(val),
	}
}

func KVInt64(k string, val int64) TraceKV {
	return TraceKV{
		key:   k,
		value: ValueInt64(val),
	}
}

func KVString(k string, val string) TraceKV {
	return TraceKV{
		key:   k,
		value: ValueString(val),
	}
}

func ValueBool(val bool) TraceValue {
	return TraceValue{
		t:         valueTypeBool,
		boolValue: val,
	}
}

func ValueString(val string) TraceValue {
	return TraceValue{
		t:        valueTypeStr,
		strValue: val,
	}
}

func ValueInt64(val int64) TraceValue {
	return TraceValue{
		t:        valueTypeInt64,
		intValue: val,
	}
}

func ValueUInt64(val uint64) TraceValue {
	return TraceValue{
		t:         valueTypeInt64,
		uintValue: val,
	}
}

func (v TraceValue) Value() any {
	switch v.t {
	case valueTypeBool:
		return v.boolValue
	case valueTypeStr:
		return v.strValue
	case valueTypeInt64:
		return v.intValue
	case valueTypeUInt64:
		return v.uintValue
	}
	return "unknown kv type"
}

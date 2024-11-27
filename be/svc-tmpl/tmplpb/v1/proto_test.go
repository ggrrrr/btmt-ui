package tmplpbv1_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func TestAsd(t *testing.T) {

	val := &tmplpbv1.TemplateData{}

	fmt.Printf("%#vv", val.ProtoReflect().Type().Descriptor())
	// lisText :=

	// proto.Message

}

func TestDataJSON(t *testing.T) {

	structValue, err := structpb.NewStruct(map[string]any{
		"struct_key_1": "struct_value_1",
		"struct_key_2": []any{123123, "key2_val2", 3123123},
	})
	require.NoError(t, err)

	testData1 := &tmplpbv1.TemplateData{
		Items: map[string]string{
			"key1": "value1",
		},
		Data: structValue,
	}

	bytes, err := json.Marshal(testData1)
	require.NoError(t, err)

	fmt.Printf("---JSON\n%v\n---JSON--\n", string(bytes))

}

func TestDataProto(t *testing.T) {

	structValue1, err := structpb.NewStruct(map[string]any{
		"struct_key_1": "struct_value_1",
		"struct_key_2": []any{123123, "key2_val2", 3123123},
	})
	require.NoError(t, err)

	testData1 := &tmplpbv1.TemplateData{
		Items: map[string]string{
			"key1": "value1",
		},
		Data: structValue1,
	}
	fmt.Printf("---ORIGIN\n%#v\n---ORIGIN--\n\n", testData1)
	fmt.Printf("---ORIGIN\n%+v\n---ORIGIN--\n\n", testData1.Data.AsMap())

	testDataBytes1, err := proto.Marshal(testData1)
	require.NoError(t, err)
	fmt.Printf("---PROTO\n%v\n---PROTO--\n", string(testDataBytes1))

	var actual1 tmplpbv1.TemplateData

	err = proto.Unmarshal(testDataBytes1, &actual1)
	require.NoError(t, err)

	fmt.Printf("---FROM PROTO[actual1]\n%#v\n---FROM PROTO--\n", &actual1)
	fmt.Printf("\t---FROM PROTO[Data]\n%#v\n---FROM PROTO--\n\n", actual1.Data)

	// anyValue2, err := actual1.Data
	require.NoError(t, err)
	// fmt.Printf("---FROM PROTO[anyValue2]\n%#v\n---FROM PROTO--\n", anyValue2)

	// fmt.Printf("\t---FROM PROTO[Data1]\n%#v\n---FROM PROTO--\n\n", actual1.Data1.AsMap())
}

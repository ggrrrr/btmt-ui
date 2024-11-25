package tmplpb_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb"
)

func TestDataJSON(t *testing.T) {
	listItems := &tmplpb.ListText{
		List: []string{"item1"},
	}

	anyValue, err := anypb.New(listItems)
	require.NoError(t, err)

	testData1 := &tmplpb.TemplateData{
		Items: map[string]string{
			"key1": "value1",
		},
		Data: anyValue,
	}

	bytes, err := json.Marshal(testData1)
	require.NoError(t, err)

	fmt.Printf("---JSON\n%v\n---JSON--\n", string(bytes))

}

func TestDataProto(t *testing.T) {

	listItems := &tmplpb.ListText{
		List: []string{"item1"},
	}

	anyValue, err := anypb.New(listItems)
	require.NoError(t, err)

	structValue1, err := structpb.NewStruct(map[string]any{
		"struct_key_1": "struct_value_1",
		"struct_key_2": []any{123123, "key2_val2", 3123123},
	})
	require.NoError(t, err)

	testData1 := &tmplpb.TemplateData{
		Items: map[string]string{
			"key1": "value1",
		},
		Data:  anyValue,
		Data1: structValue1,
	}
	fmt.Printf("---ORIGIN\n%#v\n---ORIGIN--\n\n", testData1)
	fmt.Printf("---ORIGIN\n%#v\n---ORIGIN--\n\n", testData1.Data1)
	fmt.Printf("---ORIGIN\n%+v\n---ORIGIN--\n\n", testData1.Data1.AsMap())

	testDataBytes1, err := proto.Marshal(testData1)
	require.NoError(t, err)
	fmt.Printf("---PROTO\n%v\n---PROTO--\n", string(testDataBytes1))

	var actual1 tmplpb.TemplateData

	err = proto.Unmarshal(testDataBytes1, &actual1)
	require.NoError(t, err)

	fmt.Printf("---FROM PROTO[actual1]\n%#v\n---FROM PROTO--\n", actual1)
	fmt.Printf("\t---FROM PROTO[Data]\n%#v\n---FROM PROTO--\n\n", actual1.Data)

	anyValue2, err := actual1.Data.UnmarshalNew()
	require.NoError(t, err)
	fmt.Printf("---FROM PROTO[anyValue2]\n%#v\n---FROM PROTO--\n", anyValue2)

	fmt.Printf("\t---FROM PROTO[Data1]\n%#v\n---FROM PROTO--\n\n", actual1.Data1.AsMap())
}

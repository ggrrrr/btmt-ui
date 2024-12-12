package msgbus

import (
	"fmt"
	"testing"

	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
)

func TestCmd(t *testing.T) {

	test := &msgbusv1.TestCommand{}

	fmt.Printf("%#v \n", test.ProtoReflect())
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor())
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().Name())
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().ParentFile())

}

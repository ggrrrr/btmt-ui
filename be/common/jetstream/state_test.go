package jetstream

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/state"
)

func TestKVSetter(t *testing.T) {
	ctx := context.Background()

	objectType, err := state.ParseEntityType("asdasd")
	require.NoError(t, err)

	newKey1 := "6cd9fc27-6c16-43be-a568-cc583d105212"
	newKey2 := "b89315ea-fc05-49a7-b997-0902bfdbfc7a"
	o1 := state.NewEntity{
		Key:   newKey1,
		Value: []byte("value 11"),
	}
	o2 := state.NewEntity{
		Key:   newKey2,
		Value: []byte("value 22"),
	}

	n, err := ConnectForTest()
	require.NoError(t, err)
	defer n.js.DeleteKeyValue(ctx, objectType.String())
	// n.js.kv

	testKVSetter, err := NewStateStore(ctx, n, objectType)
	require.NoError(t, err)

	rev1, err := testKVSetter.Push(ctx, o1)
	require.NoError(t, err)
	fmt.Println("rev1", rev1)

	rev2, err := testKVSetter.Push(ctx, o2)
	require.NoError(t, err)

	fmt.Println("rev2", rev2)

	testKVGetter, err := NewStateStore(ctx, n, objectType)
	require.NoError(t, err)

	var actual1, actual2 state.EntityState

	actual1, err = testKVGetter.Fetch(ctx, newKey1)
	require.NoError(t, err)

	fmt.Printf("got1[%v]: %#v \n", newKey1, rev1)

	state.TestObject(t, state.EntityState{
		Revision:  rev1,
		Key:       newKey1,
		Value:     o1.Value,
		CreatedAt: time.Now(),
	}, actual1, time.Duration(400*time.Millisecond))

	actual2, err = testKVGetter.Fetch(ctx, newKey2)
	require.NoError(t, err)

	fmt.Printf("got1[%v]: %#v \n", newKey2, actual2)

	state.TestObject(t, state.EntityState{
		Revision:  rev2,
		Key:       newKey2,
		Value:     o2.Value,
		CreatedAt: time.Now(),
	}, actual2, time.Duration(400*time.Millisecond))

	// // rev3, err = testKVGetter.GetProto(ctx, newKey2, got3)
	// // require.Error(t, err)

	// _, err = testKVGetter.History(ctx, newKey1)
	// require.NoError(t, err)

}

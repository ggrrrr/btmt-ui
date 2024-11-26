package repo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

func Test_objectId(t *testing.T) {
	uuid1 := primitive.NewObjectID()
	fmt.Printf("uuid: %v \n", uuid1)
	hex := uuid1.Hex()
	fmt.Printf("uuid: %x %v\n", hex, len(hex))

	uuid1p, err := mgo.ConvertFromId(hex)
	assert.NoError(t, err)
	assert.Equal(t, uuid1, uuid1p)
	fmt.Printf("uuid : %s \n", uuid1p.String())
	hex = uuid1p.Hex()
	fmt.Printf("uuid: %s %v\n", hex, len(hex))

	map1 := map[string]string{"key1": "val1", "key2": "val2"}
	slice1 := toSlice(map1)
	map1out := toMap(slice1)
	fmt.Printf("%v %v\n", map1, slice1)
	fmt.Printf("%v %v", map1out, slice1)
	assert.Equal(t, map1, map1out)

}

func Test_FromPerson(t *testing.T) {

	id1 := primitive.NewObjectID()

	p1 := &peoplepbv1.Person{
		Id:        id1.Hex(),
		IdNumbers: map[string]string{"pin": "pin1"},
		Emails:    map[string]string{"default": "asd@asd"},
		Name:      "asd",
		FullName:  "ewrcxf asd",
		Dob: &peoplepbv1.Dob{
			Year:  2001,
			Month: 3,
			Day:   13,
		},
		Gender:    "m",
		Phones:    map[string]string{"asd": "asdasd"},
		Labels:    []string{"red"},
		Attr:      map[string]string{"some": "1"},
		CreatedAt: timestamppb.Now(),
	}

	out1, err := fromPerson(p1)
	assert.Equal(t, out1.Emails, []string{"default:asd@asd"})
	require.NoError(t, err)
	p11 := out1.toProto()
	TestPerson(t, p11, p1, time.Duration(300*time.Millisecond))
	// assert.Equal(t, p1.DateOfBirth, p11.DateOfBirth)
}

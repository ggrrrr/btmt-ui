package repo

import (
	"fmt"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func cage(asd *int) {
	age1 := 10
	asd = &age1
}

func TestAsd(t *testing.T) {

	type Age struct {
		age *int
	}

	age := time.Now().Year() - 1978

	asd := Age{age: new(int)}
	t.Logf("%v \n", asd.age)
	cage(asd.age)
	t.Logf("%v \n", asd.age)
	asd.age = &age
	t.Logf("%v \n", *asd.age)

}

func Test_objectId(t *testing.T) {
	uuid1 := primitive.NewObjectID()
	fmt.Printf("uiud: %v \n", uuid1)
	hex := uuid1.Hex()
	fmt.Printf("uiud: %x %v\n", hex, len(hex))

	uuid1p, err := convertPersonId(hex)
	assert.NoError(t, err)
	assert.Equal(t, uuid1, uuid1p)
	fmt.Printf("uiud : %s \n", uuid1p.String())
	hex = uuid1p.Hex()
	fmt.Printf("uiud: %s %v\n", hex, len(hex))

	map1 := map[string]string{"key1": "val1", "key2": "val2"}
	slice1 := toSlice(map1)
	map1out := toMap(slice1)
	fmt.Printf("%v %v\n", map1, slice1)
	fmt.Printf("%v %v", map1out, slice1)
	assert.Equal(t, map1, map1out)

}

func Test_FromPerson(t *testing.T) {

	id1 := primitive.NewObjectID()

	p1 := ddd.Person{
		Id:        id1.Hex(),
		IdNumbers: map[string]string{"pin": "pin1"},
		Emails:    map[string]string{"default": "asd@asd"},
		Name:      "asd",
		FullName:  "ewrcxf asd",
		DOB: &ddd.Dob{
			Year:  2001,
			Month: 3,
			Day:   13,
		},
		Gender:      "m",
		Phones:      map[string]string{"asd": "asdasd"},
		Labels:      []string{"red"},
		Attr:        map[string]string{"some": "1"},
		CreatedTime: time.Now(),
	}

	out1, err := fromPerson(&p1)
	assert.Equal(t, out1.Emails, []string{"default:asd@asd"})
	require.NoError(t, err)
	p11 := out1.toPerson()
	age := (time.Now().Year() - p1.DOB.Year)
	p1.Age = &age
	TestPerson(t, p11, p1, 10)
	// assert.Equal(t, p1.DateOfBirth, p11.DateOfBirth)
}

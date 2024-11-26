package peoplepbv1

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
)

func Test_FromPerson(t *testing.T) {
	to := &ddd.Person{
		IdNumbers:   map[string]string{"egn": "egn1"},
		Emails:      map[string]string{"default": "asd@asd123"},
		LoginEmail:  "asd@asd",
		Name:        "vesko",
		Phones:      map[string]string{"mobile": "0889430425"},
		CreatedTime: time.Now(),
	}
	pbPerson := FromPerson(to)
	t.Logf("%+v \n", to)
	t.Logf("%+v \n", pbPerson)

	b, err := json.Marshal(pbPerson)
	require.NoError(t, err)
	t.Logf("%s ", string(b))
}

func TestToPerson(t *testing.T) {
	req := `{"data":{"dob":{"year":1978,"day":22,"month":2},"login_email":"asd@asd","id_numbers":{"egn":"myegn"},"emails":{"default":"asd@asd123"},"name":"vesko","phones":{"mobile":"0889430425"}}}`
	var from SaveRequest
	err := json.NewDecoder(bytes.NewReader([]byte(req))).Decode(&from)
	require.NoError(t, err)

	to := &ddd.Person{
		IdNumbers: map[string]string{"egn": "myegn"},
		DOB: &ddd.Dob{
			Year:  1978,
			Month: 2,
			Day:   22,
		},
		Emails:      map[string]string{"default": "asd@asd123"},
		LoginEmail:  "asd@asd",
		Name:        "vesko",
		Phones:      map[string]string{"mobile": "0889430425"},
		CreatedTime: from.ToPerson().CreatedTime,
	}

	assert.Equal(t, from.ToPerson(), to)
}

func TestToMap(t *testing.T) {
	empty1 := ListRequest{Filters: map[string]*ListText{}}
	map1 := empty1.ToFilter()
	assert.Equal(t, map1, map[string][]string{})

	empty2 := ListRequest{Filters: map[string]*ListText{
		"texts": {List: []string{"mytext"}},
	}}
	map2 := empty2.ToFilter()
	assert.Equal(t, map2, map[string][]string{"texts": {"mytext"}})

}

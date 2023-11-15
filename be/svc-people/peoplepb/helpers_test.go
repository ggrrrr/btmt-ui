package peoplepb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

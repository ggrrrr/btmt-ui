package repo

import (
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/stretchr/testify/assert"
)

func TestPerson(t *testing.T, got *ddd.Person, want *ddd.Person, duration int) {
	if duration > 0 {
		assert.WithinDuration(t, want.CreatedTime, got.CreatedTime, 100+time.Millisecond)
	}
	// assert.WithinDuration(t, want.DateOfBirth, got.DateOfBirth, 100+time.Millisecond)
	// got.DateOfBirth = time.Time{}
	got.CreatedTime = time.Time{}
	// want.DateOfBirth = time.Time{}
	want.CreatedTime = time.Time{}
	assert.Equal(t, want, got)
}

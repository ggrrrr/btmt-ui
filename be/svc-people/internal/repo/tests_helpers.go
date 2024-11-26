package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

func TestPerson(t *testing.T, expected *peoplepbv1.Person, actual *peoplepbv1.Person, duration time.Duration) {
	if duration > 0 {
		assert.WithinDurationf(t, expected.CreatedAt.AsTime(), actual.CreatedAt.AsTime(), duration, "CreatedAt expected: %v actual: %v", expected.CreatedAt.AsTime(), actual.CreatedAt.AsTime())
		assert.WithinDurationf(t, expected.UpdatedAt.AsTime(), actual.UpdatedAt.AsTime(), duration, "UpdatedAt expected: %v actual: %v", expected.UpdatedAt.AsTime(), actual.UpdatedAt.AsTime())
	}
	// assert.WithinDuration(t, want.DateOfBirth, got.DateOfBirth, 100+time.Millisecond)
	// got.DateOfBirth = time.Time{}
	expected.CreatedAt = actual.CreatedAt
	expected.UpdatedAt = actual.UpdatedAt
	// want.DateOfBirth = time.Time{}
	// want.CreatedAt = time.Time{}
	// assert.Equal(t, want, got)
	// require.Equal(t, 1, 2, "")

	// fmt.Printf("\n\n%#v %#v\n\n", exp, act)

	if !proto.Equal(expected, actual) {
		require.Truef(t, false, "expected:\n-> %+v\n-> %+v\n", expected, actual)

	}
}

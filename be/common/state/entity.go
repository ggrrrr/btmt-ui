package state

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type (
	EntityType struct {
		nameVersion string
	}

	EntityState struct {
		Revision  uint64
		Key       string
		Value     []byte
		CreatedAt time.Time
	}

	NewEntity struct {
		Key   string
		Value []byte
	}
)

func EntityTypeFromProto(msg proto.Message) EntityType {
	return EntityType{
		nameVersion: formatBucketName(string(msg.ProtoReflect().Descriptor().FullName())),
	}
}

func ParseEntityType(name string) (EntityType, error) {
	return EntityType{
		nameVersion: formatBucketName(name),
	}, nil
}

func MustParseEntityType(name string) EntityType {
	return EntityType{
		nameVersion: formatBucketName(name),
	}
}

func (o EntityType) String() string {
	return o.nameVersion
}

func formatBucketName(name string) string {
	return strings.ReplaceAll(string(name), ".", "_")
}

func TestObject(t *testing.T, exp EntityState, actual EntityState, delta time.Duration) bool {
	if delta > 0 {
		require.WithinDurationf(t, exp.CreatedAt, actual.CreatedAt, delta, "exp: %v, actual: %v", exp.CreatedAt, actual.CreatedAt)
	}
	exp.CreatedAt = actual.CreatedAt
	require.Equal(t, exp, actual)
	return true
}

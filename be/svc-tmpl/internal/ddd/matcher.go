package ddd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MatchTemplate(t *testing.T, expected Template, actual Template) bool {
	delta := 200 * time.Millisecond
	if !assert.WithinDurationf(t, expected.UpdatedAt, actual.UpdatedAt, delta, "UpdatedAt: expected: %v actual:%v", expected.UpdatedAt, actual.UpdatedAt) {
		return false
	}
	if !assert.WithinDurationf(t, expected.CreatedAt, actual.CreatedAt, delta, "CreatedAt: expected: %v actual:%v", expected.CreatedAt, actual.CreatedAt) {
		return false
	}

	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt

	return assert.Equal(t, expected, actual)

}

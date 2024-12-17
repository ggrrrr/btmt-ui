package tmplpbv1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MatchTemplateUpdate(t *testing.T, ts time.Time, images []string, expected *TemplateUpdate, actual *Template) {
	delta := 200 * time.Millisecond
	if !assert.WithinDurationf(t, ts, actual.UpdatedAt.AsTime(), delta, "UpdatedAt: expected: %v actual:%v", ts, actual.UpdatedAt) {
		return
	}
	if !assert.WithinDurationf(t, ts, actual.CreatedAt.AsTime(), delta, "CreatedAt: expected: %v actual:%v", ts, actual.CreatedAt) {
		return
	}

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.Labels, actual.Labels)
	assert.Equal(t, expected.Body, actual.Body)
	assert.Equal(t, images, actual.Images)

}

func MatchTemplate(t *testing.T, ts time.Time, expected *Template, actual *Template) {
	delta := 200 * time.Millisecond
	if !assert.WithinDurationf(t, ts, actual.UpdatedAt.AsTime(), delta, "UpdatedAt: expected: %v actual:%v", ts, actual.UpdatedAt) {
		return
	}
	if !assert.WithinDurationf(t, ts, actual.CreatedAt.AsTime(), delta, "CreatedAt: expected: %v actual:%v", ts, actual.CreatedAt) {
		return
	}

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.Labels, actual.Labels)
	assert.Equal(t, expected.Body, actual.Body)
	assert.Equal(t, expected.Images, actual.Images)

}

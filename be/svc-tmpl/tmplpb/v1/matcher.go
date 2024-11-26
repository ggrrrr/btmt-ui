package tmplpbv1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

func MatchTemplate(t *testing.T, expected *Template, actual *Template) bool {
	delta := 200 * time.Millisecond
	if !assert.WithinDurationf(t, expected.UpdatedAt.AsTime(), actual.UpdatedAt.AsTime(), delta, "UpdatedAt: expected: %v actual:%v", expected.UpdatedAt, actual.UpdatedAt) {
		return false
	}
	if !assert.WithinDurationf(t, expected.CreatedAt.AsTime(), actual.CreatedAt.AsTime(), delta, "CreatedAt: expected: %v actual:%v", expected.CreatedAt, actual.CreatedAt) {
		return false
	}

	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt

	ok := proto.Equal(expected, actual)
	require.Truef(t, ok, "expected: %#v, actual: %#v ")
	return ok
}

var _ (logger.TraceDataExtractor) = (*Template)(nil)

func (t *Template) Extract() logger.TraceData {
	out := logger.TraceData{
		"tmpl.Name":        logger.TraceValueString(t.Name),
		"tmpl.ContentType": logger.TraceValueString(t.ContentType),
	}
	if t.Id != "" {
		out["tmpl.Id"] = logger.TraceValueString(t.Id)
	}
	return out
}

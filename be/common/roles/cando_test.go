package roles

import (
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	test     string
	testFunc func(tt *testing.T)
}

var err error

func Test_Ok(t *testing.T) {
	testCanDo := NewAppPolices()

	tests := []testCase{
		{
			test: "ErrAuthUnauthenticated",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo("", AuthInfo{})
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
			},
		},
		{
			test: "ErrForbidden",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo("", AuthInfo{User: "asd"})
				assert.ErrorIs(tt, err, app.ErrForbidden)
			},
		},
		{
			test: "ok admin one role",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo("", AuthInfo{User: "asd", Roles: []RoleName{"admin"}})
				assert.NoError(t, err)
			},
		},
		{
			test: "ok admin more roles",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo("", AuthInfo{User: "asd", Roles: []RoleName{"asd", "admin"}})
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}
}

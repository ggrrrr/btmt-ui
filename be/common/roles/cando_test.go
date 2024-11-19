package roles

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

var testTenant string = "test-host"

func Test_CanDo(t *testing.T) {
	var err error

	testCanDo := NewAppPolices()

	tests := []struct {
		test     string
		testFunc func(tt *testing.T)
	}{
		{
			test: "ok custom role ok",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(testTenant, "my-new-role", AuthInfo{Realm: testTenant, Subject: "asd", Roles: []string{"my-new-role"}})
				assert.NoError(tt, err)
			},
		},
		{
			test: "ErrAuthUnauthenticated",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(testTenant, "", AuthInfo{})
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
			},
		},
		{
			test: "ErrForbidden",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(testTenant, "", AuthInfo{Subject: "asd"})
				assert.ErrorIs(tt, err, app.ErrForbidden)
			},
		},
		{
			test: "ok admin one role",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(testTenant, "", AuthInfo{Realm: testTenant, Subject: "asd", Roles: []string{"admin"}})
				assert.NoError(tt, err)
			},
		},
		{
			test: "ok admin more roles",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(testTenant, "", AuthInfo{Realm: testTenant, Subject: "asd", Roles: []string{"asd", "admin"}})
				assert.NoError(tt, err)
			},
		},
		{
			test: "ok system ErrForbidden",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(SystemRealm, "", AuthInfo{Realm: testTenant, Subject: "asd", Roles: []string{"asd", "adasdmin"}})
				assert.ErrorIs(tt, err, app.ErrForbidden)
			},
		},
		{
			test: "ok system ok",
			testFunc: func(tt *testing.T) {
				err = testCanDo.CanDo(SystemRealm, "", AuthInfo{Realm: testTenant, Subject: "asd", SystemRoles: []string{"asd", "admin"}})
				assert.NoError(tt, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}
}

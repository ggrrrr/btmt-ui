package app

import (
	"context"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/mongodb"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	testCase struct {
		test     string
		testFunc func(tt *testing.T)
	}
)

func TestSave(t *testing.T) {
	rootCtx := context.Background()
	ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("mock", roles.Device{}))
	ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})
	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	cfg := mongodb.Config{
		TTL:        10 * time.Second,
		Collection: "app.TestList",
		User:       "admin",
		Passwd:     "pass",
		Database:   "people",
		Url:        "mongodb://localhost:27017/",
		Debug:      "console",
	}
	testDb, err := mongodb.New(rootCtx, cfg)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)

	testRepo := repo.New(cfg.Collection, testDb)

	testApp, err := New(
		WithPeopleRepo(testRepo),
		WithAppPolicies(roles.NewAppPolices()),
	)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "no access",
			testFunc: func(tt *testing.T) {
				_, err = testApp.GetById(rootCtx, "asd")
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
				err := testApp.Save(ctxNormal, &ddd.Person{
					PIN:    "123123",
					Email:  "new@asd",
					Name:   "new",
					Phones: map[string]string{"mobile": "123123"},
				})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				_, err = testApp.GetById(ctxNormal, "asd")
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Update(ctxNormal, &ddd.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Update(ctxNormal, &ddd.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Save(ctxNormal, &ddd.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
			},
		},
		{
			test: "save",
			testFunc: func(tt *testing.T) {
				p1 := &ddd.Person{
					PIN:      "pin1",
					Name:     "name 1",
					Email:    "email 1",
					FullName: "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				err := testApp.Save(ctxAdmin, p1)
				assert.NoError(tt, err)
				p2, err := testApp.GetById(ctxAdmin, p1.Id)
				assert.NoError(tt, err)
				repo.TestPerson(tt, *p2, *p1, 10)
			},
		},
		{
			test: "forbiden",
			testFunc: func(tt *testing.T) {
				p1 := &ddd.Person{
					PIN:      "pin1",
					Name:     "name 1",
					Email:    "email 1",
					FullName: "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				err := testApp.Save(ctxNormal, p1)
				assert.Error(tt, err)
				assert.ErrorIs(tt, err, app.ErrForbidden)
				// p2, err := testApp.GetById(ctxAdmin, p1.Id)
				// assert.NoError(tt, err)
				// repo.TestPerson(tt, *p2, *p1, 10)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}
}

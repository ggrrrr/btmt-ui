package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
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
	ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemTenant, "mock", roles.Device{}))
	ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})

	cfg := mgo.MgoTestCfg()
	testDb, err := mgo.New(rootCtx, cfg)
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
			test: "no accessa",
			testFunc: func(tt *testing.T) {
				_, err = testApp.GetById(rootCtx, "asd")
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
				err := testApp.Save(ctxNormal, &ddd.Person{
					IdNumbers:  map[string]string{"pin": "pin1"},
					LoginEmail: "new@asd",
					Emails:     map[string]string{"g": "new@asd"},
					Name:       "new",
					Phones:     map[string]string{"mobile": "123123"},
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
					IdNumbers:  map[string]string{"pin": "pin1"},
					Name:       "name 1",
					LoginEmail: "email 1",
					Emails:     map[string]string{"g": "asd@asd"},
					FullName:   "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				ts := time.Now()

				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				p2, err := testApp.GetById(ctxAdmin, p1.Id)
				require.NoError(tt, err)
				assert.WithinDuration(tt, ts, p1.CreatedTime, 1+time.Second)
				tt.Logf("p1: %v \n", p1)
				repo.TestPerson(tt, *p2, *p1, 10)
				list, err := testApp.List(ctxAdmin, nil)
				require.NoError(tt, err)
				require.True(tt, len(list) > 0)
				tt.Logf("list: %v \n", len(list))
				// require.True(tt, list[0].CreatedTime == 1)
				tt.Logf("list: %v \n", list[0].CreatedTime)

			},
		},
		{
			test: "forbidden",
			testFunc: func(tt *testing.T) {
				p1 := &ddd.Person{
					IdNumbers:  map[string]string{"pin": "pin1"},
					Name:       "name 1",
					LoginEmail: "email 1",
					Emails:     map[string]string{"g1": "asd@asd1"},
					FullName:   "full name 1",
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
		{
			test: "pin validator",
			testFunc: func(tt *testing.T) {
				p1 := &ddd.Person{
					IdNumbers:  map[string]string{"EGN": os.Getenv("PIN2")},
					Name:       "name 1",
					LoginEmail: "email 1",
					Emails:     map[string]string{"g1": "asd@asd1"},
					FullName:   "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				assert.Equal(tt, p1.DOB, &ddd.Dob{Year: 1978, Month: 2, Day: 13})
				tt.Logf("%+v \n", p1)
				age := time.Now().Year() - p1.DOB.Year
				p1.Age = &age

				p2, err := testApp.GetById(ctxAdmin, p1.Id)
				require.NoError(tt, err)
				assert.Equal(tt, p2.Gender, "male")
				assert.Equal(tt, p1.Gender, "male")
				repo.TestPerson(tt, *p2, *p1, 10)
			},
		},
		{
			test: "update",
			testFunc: func(tt *testing.T) {
				p1 := &ddd.Person{}
				p2 := &ddd.Person{
					IdNumbers:  map[string]string{"pin": "pin1"},
					LoginEmail: "loginemail",
					Name:       "name1",
					FullName:   "full name",
					DOB: &ddd.Dob{
						Year:  1978,
						Month: 2,
						Day:   13,
					},
					Gender: "male",
					Emails: map[string]string{"main": "main@mail"},
					Phones: map[string]string{"main": "phone1"},
					Labels: []string{"atttre"},
					Attr:   map[string]string{"attr": "atttrvalue"},
				}
				p2.CreatedTime = time.Now()
				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				tt.Logf("%+v \n", p1)
				p2.Id = p1.Id
				err = testApp.Update(ctxAdmin, p2)
				require.NoError(tt, err)
				p3, err := testApp.GetById(ctxAdmin, p2.Id)
				p3.Age = nil
				require.NoError(tt, err)
				repo.TestPerson(tt, *p3, *p2, 5)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}
}

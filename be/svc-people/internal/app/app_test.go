package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	testCase struct {
		test     string
		testFunc func(tt *testing.T)
	}
)

func TestSave(t *testing.T) {

	createAdDurtion := 200 * time.Millisecond

	rootCtx := context.Background()
	ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemRealm, "mock", app.Device{}))
	ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{Subject: "some"})

	cfg := mgo.MgoTestCfg("test-people")
	testDb, err := mgo.New(rootCtx, cfg)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)

	testRepo := repo.New(cfg.Collection, testDb)

	mockState := &state.MockStore{}

	testApp, err := New(
		WithPeopleRepo(testRepo),
		WithAppPolicies(roles.NewAppPolices()),
		WithStateStore(mockState),
	)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "no accessa",
			testFunc: func(tt *testing.T) {
				_, err = testApp.GetById(rootCtx, "asd")
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
				mockState.On("Push", mock.Anything).Return(uint64(1), nil)
				err := testApp.Save(ctxNormal, &peoplepbv1.Person{
					IdNumbers:    map[string]string{"pin": "pin1"},
					PrimaryEmail: "new@asd",
					Emails:       map[string]string{"g": "new@asd"},
					Name:         "new",
					Phones:       map[string]string{"mobile": "123123"},
				})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				_, err = testApp.GetById(ctxNormal, "asd")
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Update(ctxNormal, &peoplepbv1.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Update(ctxNormal, &peoplepbv1.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
				err = testApp.Save(ctxNormal, &peoplepbv1.Person{})
				assert.ErrorIs(tt, err, app.ErrForbidden)
			},
		},
		{
			test: "save",
			testFunc: func(tt *testing.T) {
				p1 := &peoplepbv1.Person{
					IdNumbers:    map[string]string{"pin": "pin1"},
					Name:         "name 1",
					PrimaryEmail: "email 1",
					Emails:       map[string]string{"g": "asd@asd"},
					FullName:     "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				ts := time.Now()

				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				p2, err := testApp.GetById(ctxAdmin, p1.Id)
				require.NoError(tt, err)
				assert.WithinDuration(tt, ts, p1.CreatedAt.AsTime(), createAdDurtion)
				tt.Logf("p1: %v \n", p1)
				repo.TestPerson(tt, p2, p1, createAdDurtion)
				list, err := testApp.List(ctxAdmin, nil)
				require.NoError(tt, err)
				require.True(tt, len(list) > 0)
				tt.Logf("list: %v \n", len(list))
				// require.True(tt, list[0].CreatedTime == 1)
				tt.Logf("list: %v \n", list[0].CreatedAt)

			},
		},
		{
			test: "forbidden",
			testFunc: func(tt *testing.T) {
				p1 := &peoplepbv1.Person{
					IdNumbers:    map[string]string{"pin": "pin1"},
					Name:         "name 1",
					PrimaryEmail: "email 1",
					Emails:       map[string]string{"g1": "asd@asd1"},
					FullName:     "full name 1",
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
				p1 := &peoplepbv1.Person{
					IdNumbers:    map[string]string{"EGN": os.Getenv("PIN2")},
					Name:         "name 1",
					PrimaryEmail: "email 1",
					Emails:       map[string]string{"g1": "asd@asd1"},
					FullName:     "full name 1",
					// Phones:   map[string]string{"mobile": "phone1"},
					// Attr:     make(map[string]string),
					// Labels: []string{},
				}
				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				assert.Equal(tt, p1.Dob.Day, uint32(13))
				assert.Equal(tt, p1.Dob.Month, uint32(2))
				assert.Equal(tt, p1.Dob.Year, uint32(1978))
				tt.Logf("%+v \n", p1)

				p2, err := testApp.GetById(ctxAdmin, p1.Id)
				require.NoError(tt, err)
				assert.Equal(tt, p2.Gender, "male")
				assert.Equal(tt, p1.Gender, "male")
				repo.TestPerson(tt, p2, p1, createAdDurtion)
			},
		},
		{
			test: "update",
			testFunc: func(tt *testing.T) {
				p1 := &peoplepbv1.Person{}
				expected := &peoplepbv1.Person{
					IdNumbers:    map[string]string{"pin": "pin1"},
					PrimaryEmail: "loginemail",
					Name:         "name1",
					FullName:     "full name",
					Dob: &peoplepbv1.Dob{
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
				err := testApp.Save(ctxAdmin, p1)
				require.NoError(tt, err)
				expected.Id = p1.Id
				expected.CreatedAt = p1.CreatedAt
				err = testApp.Update(ctxAdmin, expected)
				require.NoError(tt, err)
				actual, err := testApp.GetById(ctxAdmin, expected.Id)
				require.NoError(tt, err)
				repo.TestPerson(tt, expected, actual, createAdDurtion)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}
}

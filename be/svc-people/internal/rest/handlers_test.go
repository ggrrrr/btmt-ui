package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonApp "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

var (
	collection string = "test-people"
)

func Test_Save(t *testing.T) {
	rootCtx := context.Background()

	// ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("mock", roles.Device{}))
	// ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})
	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	testDb, err := mgo.ConnectForTest(collection)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)
	mockState := &state.MockStore{}
	testRepo := repo.New(collection, testDb)
	app, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
		app.WithStateStore(mockState),
	)
	require.NoError(t, err)

	testServer := server{tracer: tracer.Tracer(otelScope), app: app}
	w := httptest.NewRecorder()
	reqStr := `{"data":{"email":"asd@asd123","name":"vesko","phones":{"mobile":"0889430425"}}}`
	httpReq := httptest.NewRequest(http.MethodPost, "/greet?name=john", strings.NewReader(reqStr))
	httpReq = httpReq.WithContext(roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemRealm, "asd", commonApp.Device{})))
	mockState.On("Push", mock.Anything).Return(uint64(1), nil)
	testServer.Save(w, httpReq)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	responeBytes, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	fmt.Printf("responeBytes %v", string(responeBytes))
	res := &peoplepbv1.SaveResponse{}
	err = json.Unmarshal(responeBytes, res)
	require.NoError(t, err)

	require.True(t, res.Payload.Id != "")
	actual, err := testRepo.GetById(rootCtx, res.Payload.Id)
	require.NoError(t, err)

	fmt.Printf("actual %+v", actual)
	expected := &peoplepbv1.Person{
		Id:        res.Payload.Id,
		Name:      "vesko",
		Phones:    map[string]string{"mobile": "0889430425"},
		CreatedAt: timestamppb.Now(),
	}
	repo.TestPerson(t, expected, actual, time.Duration(400*time.Millisecond))

}

func Test_List(t *testing.T) {
	rootCtx := context.Background()

	testDb, err := mgo.ConnectForTest(collection)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)

	testRepo := repo.New(collection, testDb)
	app, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
		app.WithStateStore(&state.MockStore{}),
	)
	require.NoError(t, err)

	testServer := server{tracer: tracer.Tracer(otelScope), app: app}
	w := httptest.NewRecorder()
	reqStr := `{}`
	httpReq := httptest.NewRequest(http.MethodPost, "/greet?name=john", strings.NewReader(reqStr))
	httpReq = httpReq.WithContext(roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemRealm, "asd", commonApp.Device{})))
	testServer.List(w, httpReq)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	asd, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	fmt.Printf("JSON %v\n", string(asd))

}

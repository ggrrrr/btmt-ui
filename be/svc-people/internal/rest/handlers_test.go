package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Save(t *testing.T) {
	rootCtx := context.Background()
	// ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("mock", roles.Device{}))
	// ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})
	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	cfg := mgo.MgoTestCfg()
	testDb, err := mgo.New(rootCtx, cfg)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)

	testRepo := repo.New(cfg.Collection, testDb)
	app, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
	)
	require.NoError(t, err)

	testServer := server{app: app}
	w := httptest.NewRecorder()
	reqStr := `{"data":{"pin":"asdads","email":"asd@asd123","name":"vesko","phones":{"mobile":"0889430425"}}}`
	httpReq := httptest.NewRequest(http.MethodPost, "/greet?name=john", strings.NewReader(reqStr))
	httpReq = httpReq.WithContext(roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemRealm, "asd", roles.Device{})))
	testServer.Save(w, httpReq)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	asd, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	fmt.Printf("asd %v", string(asd))

}

func Test_List(t *testing.T) {
	rootCtx := context.Background()
	// ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("mock", roles.Device{}))
	// ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})
	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	cfg := mgo.MgoTestCfg()
	testDb, err := mgo.New(rootCtx, cfg)
	require.NoError(t, err)
	defer testDb.Close(rootCtx)

	testRepo := repo.New(cfg.Collection, testDb)
	app, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
	)
	require.NoError(t, err)

	testServer := server{app: app}
	w := httptest.NewRecorder()
	reqStr := `{}`
	httpReq := httptest.NewRequest(http.MethodPost, "/greet?name=john", strings.NewReader(reqStr))
	httpReq = httpReq.WithContext(roles.CtxWithAuthInfo(rootCtx, roles.CreateSystemAdminUser(roles.SystemRealm, "asd", roles.Device{})))
	testServer.List(w, httpReq)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	asd, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	fmt.Printf("JSON %v\n", string(asd))

}

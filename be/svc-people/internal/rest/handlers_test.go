package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/mongodb"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Asd(t *testing.T) {
	rootCtx := context.Background()
	// ctxAdmin := roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("mock", roles.Device{}))
	// ctxNormal := roles.CtxWithAuthInfo(rootCtx, roles.AuthInfo{User: "some"})
	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	cfg := mongodb.Config{
		TTL:        10 * time.Second,
		Collection: "people",
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
	app, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
	)
	require.NoError(t, err)

	testServer := server{app: app}
	w := httptest.NewRecorder()
	reqStr := `{"data":{"pin":"asdads","email":"asd@asd123","name":"vesko","phones":{"mobile":"0889430425"}}}`
	httpReq := httptest.NewRequest(http.MethodPost, "/greet?name=john", strings.NewReader(reqStr))
	httpReq = httpReq.WithContext(roles.CtxWithAuthInfo(rootCtx, roles.CreateAdminUser("asd", roles.Device{})))
	testServer.Save(w, httpReq)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
}

package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
)

type testApp struct {
	wg *sync.WaitGroup
	t  *testing.T
}

func (a *testApp) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("testHandler %+v %+v \n", r.Method, r.URL.Path)
	body := "ok"
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	bytes, err := w.Write([]byte(body))
	assert.NoError(a.t, err)
	assert.Equal(a.t, 2, bytes)

	authInfo := roles.AuthInfoFromCtx(r.Context())
	assert.Equal(a.t, "//Go-http-client", authInfo.Device.DeviceInfo)
	assert.Equal(a.t, "[::1]", authInfo.Device.RemoteAddr)
	assert.Equal(a.t, "tosho", authInfo.Subject)
	logger.DebugCtx(r.Context()).Msg("testHandler")
	a.wg.Done()
}

func TestRest(t *testing.T) {
	app := &testApp{
		t:  t,
		wg: &sync.WaitGroup{},
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	testSystem := &System{
		cfg: config.AppConfig{
			Rest: config.RestConfig{
				Address: ":9091",
			},
			Jwt: config.JwtConfig{
				UseMock: "mock",
			},
		},
	}
	err := testSystem.initJwt()
	require.NoError(t, err)
	testSystem.initMux()
	testSystem.waiter = waiter.New()

	testSystem.Mux().Get("/test", app.testHandler)

	go func() {
		err = testSystem.WaitForWeb(ctx)
		require.NoError(t, err)
	}()

	app.wg.Add(1)
	req, err := http.NewRequest("GET", "http://localhost:9091/test", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "mock tosho")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	_ = resp.Body.Close()
	require.Equalf(t, "ok", string(respBody), "%+v", respBody)
	app.wg.Wait()
	cancelFunc()
}

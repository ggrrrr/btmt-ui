package system

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	app.wg.Add(2)
	// testSystem.Mux().Mount("/test", app.testHandler)

	resp, err := http.Get("http://localhost:9091/test")
	require.NoError(t, err)
	// require.NotNil(t, resp)
	var respBody []byte
	_, err = resp.Body.Read(respBody)
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equalf(t, "ok", string(respBody), "%+v", respBody)
	app.wg.Wait()
	cancelFunc()
}

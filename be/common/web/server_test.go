package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	t *testing.T
}

func getTestPort(s *Server) int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func getTestUrl(s *Server) string {
	return fmt.Sprintf("http://localhost:%d", getTestPort(s))
}

func (ts *testServer) testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	_, _ = w.Write([]byte(`.`))
}

func TestStartupShutdown(t *testing.T) {
	var err error
	ts := &testServer{
		t: t,
	}
	ctx := context.Background()
	testServer1, err := NewServer("", Config{
		// EndpointREST: "/rest",
		ListenAddr: ":0",
		CORS:       CORS{Origin: "origin", Headers: "CORS-Header"},
	})
	require.NoError(t, err)
	testServer1.MountHandler("/app1", http.HandlerFunc(ts.testHandler))
	go func() {
		err = testServer1.Startup()
		require.NoError(t, err)
	}()
	defer func() {
		err = testServer1.Shutdown(ctx)
		assert.NoError(t, err)
	}()

	// fmt.Println("asd")
	// time.Sleep(2 * time.Second)
	testServer1.listenReady.Wait()

	resp, err := http.Get(getTestUrl(testServer1) + "/rest/app1/assssssd")
	require.NoError(t, err)
	fmt.Printf("got: %+v \n", resp.Status)
	fmt.Printf("got: %+v \n", resp.Header)

	testServer2, err := NewServer("server2", Config{
		ListenAddr: fmt.Sprintf(":%d", getTestPort(testServer1)),
	})
	require.NoError(t, err)

	err = testServer2.Startup()
	require.Error(t, err)

	err = testServer1.Shutdown(ctx)
	require.NoError(t, err)

}

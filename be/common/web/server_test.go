package web

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	_, _ = w.Write([]byte(`.`))
}

func TestStartupShutdown(t *testing.T) {
	var err error
	ctx := context.Background()
	testServer1, err := NewServer("", Config{
		// EndpointREST: "/rest",
		ListenAddr: ":8080",
	})
	require.NoError(t, err)
	testServer1.MountHandler("/app1", http.HandlerFunc(testHandler))

	go func() {
		err = testServer1.Startup()
		require.NoError(t, err)
	}()
	defer func() {
		err = testServer1.Shutdown(ctx)
		assert.NoError(t, err)
	}()

	fmt.Println("asd")
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8080/rest/app1/assssssd")
	require.NoError(t, err)
	fmt.Printf("got: %+v \n", resp.Status)

	testServer2, err := NewServer("server2", Config{
		// EndpointREST: "/rest2",
		ListenAddr: ":8080",
	})
	require.NoError(t, err)

	err = testServer2.Startup()
	require.Error(t, err)

	err = testServer1.Shutdown(ctx)
	require.NoError(t, err)

}

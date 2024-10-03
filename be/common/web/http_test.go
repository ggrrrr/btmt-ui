package web

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func myMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Printf("processing request...\n")
		defer fmt.Printf("defer\n")

		// We use `select` to execute a piece of code depending on which
		// channel receives a message first
		<-time.After(2 * time.Second)
		// If we receive a message after 2 seconds
		// that means the request has been processed
		// We then write this as the response
		_, _ = w.Write([]byte("request processed"))
		fmt.Printf("processing request.done.\n")
		// case <-ctx.Done():
		// 	// If the request gets cancelled, log it
		// 	// to STDERR
		// 	fmt.Printf("request cancelled\n")
		// }
		// fmt.Fprint(os.Stderr, "request cancelled\n")
	})
	return mux
}

func TestHTTP(t *testing.T) {
	addr := "localhost:10801"
	webServer := &http.Server{
		Addr:    addr,
		Handler: myMux(),
	}

	go func() {
		err := webServer.ListenAndServe()
		require.NoError(t, err)
	}()
	// defer webServer.Close()
	time.Sleep(2 * time.Second)

	reqCtx, reqCancel := context.WithTimeout(context.Background(), 1*time.Second)

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, fmt.Sprintf("http://%s", addr), nil)
	require.NoError(t, err)
	defer reqCancel()

	client := &http.Client{}
	res, err := client.Do(req)
	require.Error(t, err)
	if err == nil {
		res.Body.Close()
	}

	fmt.Printf("done \n")
}

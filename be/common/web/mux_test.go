package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAppRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(r.Context(), w, 200, "", nil, "test val")
	})
	return r

}

func Test_Server(t *testing.T) {

	tests := []struct {
		name     string
		method   string
		endpoint string
		prepFn   func(t *testing.T) *Server
		code     int
		jsonBody string
		textBody string
		headers  string
	}{
		{
			name:     "ok",
			method:   "GET",
			endpoint: "/testApp/test",
			code:     200,
			jsonBody: `{"code":"200", "payload":"test val"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"server",
					Config{
						ListenAddr: ":8081",
					})
				require.NoError(t, err)

				testServer.MountHandler("/testApp", testAppRouter())

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
		{
			name:     "not found",
			method:   "GET",
			endpoint: "/test1",
			code:     404,
			jsonBody: `{"code":"404", "message":"NotFound: [GET] /test1"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":8081",
						CORS: CORS{
							Origin:  "*",
							Headers: "HeaderName",
						},
					})
				require.NoError(t, err)

				testServer.mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
					sendJSON(r.Context(), w, 200, "ok", nil, "test")
				})

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
		{
			name:     "options",
			method:   http.MethodOptions,
			endpoint: "/test",
			code:     200,
			textBody: `.`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":8081",
						CORS: CORS{
							Origin:  "*",
							Headers: "HeaderName",
						},
					})
				require.NoError(t, err)

				testServer.mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
					sendJSON(r.Context(), w, 200, "ok", nil, "test")
				})

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
		{
			name:     "panic",
			method:   http.MethodGet,
			endpoint: "/panic",
			code:     500,
			textBody: `.`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":8081",
					})
				require.NoError(t, err)
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

				testServer.mux.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
					fmt.Printf("PANIC PANIC PANIC PANIC about to panic\n")
					panic(fmt.Errorf("test panic"))
				})

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testServer := tc.prepFn(t)
			defer testServer.Shutdown(context.Background())

			httpReq, err := http.NewRequest(tc.method, "http://localhost:8081"+tc.endpoint, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(httpReq)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			resp.Body.Close()

			fmt.Printf("body %s \n", string(body))
			fmt.Printf("headers %d %#v \n", resp.StatusCode, resp.Header)

			if tc.jsonBody != "" {
				assert.JSONEq(t, tc.jsonBody, string(body))
			}
			if tc.textBody != "" {
				assert.Equal(t, tc.textBody, string(body))
			}

		})
	}

}

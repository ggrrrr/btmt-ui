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

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

func testAppRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(r.Context(), w, 200, "", nil, "test val")
	})
	return r
}

func testAuthAppRouter(t *testing.T, expected string) http.Handler {
	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		user := roles.AuthInfoFromCtx(r.Context())
		if user.Subject == "" {
			SendJSONForbidden(r.Context(), w, "forbidden")
			return
		}
		assert.Equal(t, expected, user.Subject)
		sendJSON(r.Context(), w, 200, "", nil, "test val")
	})
	return r
}

func Test_Server(t *testing.T) {

	tests := []struct {
		name       string
		method     string
		endpoint   string
		reqHeaders http.Header
		prepFn     func(t *testing.T) *Server
		code       int
		jsonBody   string
		textBody   string
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
						ListenAddr: ":0",
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
						ListenAddr: ":0",
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
						ListenAddr: ":0",
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
			textBody: `{"code":"500","message":"internal panic"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":0",
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
		{
			name:     "auth 403",
			method:   http.MethodGet,
			endpoint: "/testApp/test",
			code:     403,
			textBody: `{"code":"403","message":"forbidden"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":0",
					},
					WithVerifier(token.NewVerifierMock()),
				)
				require.NoError(t, err)
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

				testServer.MountHandler("/testApp", testAuthAppRouter(t, "asd"))

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
		{
			name:       "auth 200",
			method:     http.MethodGet,
			reqHeaders: http.Header{"Authorization": []string{"Bearer username"}},
			endpoint:   "/testApp/test",
			code:       200,
			textBody:   `{"code":"200","payload":"test val"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":0",
					},
					WithVerifier(token.NewVerifierMock()),
				)
				require.NoError(t, err)
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

				testServer.MountHandler("/testApp", testAuthAppRouter(t, "username"))

				go func() {
					err = testServer.Startup()
					require.NoError(t, err)
				}()

				return testServer
			},
		},
		{
			name:       "auth 401",
			method:     http.MethodGet,
			reqHeaders: http.Header{"Authorization": []string{"Shit username"}},
			endpoint:   "/testApp/test",
			code:       401,
			textBody:   `{"code":"401","message":"Unauthenticated"}`,
			prepFn: func(t *testing.T) *Server {
				var err error
				testServer, err := NewServer(
					"",
					Config{
						ListenAddr: ":0",
					},
					WithVerifier(token.NewVerifierMock()),
				)
				require.NoError(t, err)
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

				testServer.MountHandler("/testApp", testAuthAppRouter(t, "username"))

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
			testServer.listenReady.Wait()

			httpReq, err := http.NewRequest(tc.method, getTestUrl(testServer)+tc.endpoint, nil)
			require.NoError(t, err)

			httpReq.Header = tc.reqHeaders

			resp, err := http.DefaultClient.Do(httpReq)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			resp.Body.Close()

			fmt.Printf("body %s \n", string(body))
			fmt.Printf("headers %d %#v \n", resp.StatusCode, resp.Header)

			assert.Equal(t, tc.code, resp.StatusCode)

			if tc.jsonBody != "" {
				assert.JSONEq(t, tc.jsonBody, string(body))
			}
			if tc.textBody != "" {
				assert.Equal(t, tc.textBody, string(body))
			}

		})
	}

}

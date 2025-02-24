package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

func (s *Server) handlerRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rvr := recover()
			if rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}
				// logger.ErrorCtx(r.Context(), fmt.Errorf("recover %+v", err)).
				// 	Str("stack", string(debug.Stack())).
				// 	Msg("httpRecovery")

				// fmt.Println(string(debug.Stack()))
				fmt.Printf("\nrecover: %#v\n\n", rvr)
				logger.ErrorCtx(r.Context(), fmt.Errorf("panic")).Any("panic", rvr).Stack().Send()

				SendJSONSystemError(r.Context(), w, "internal panic", nil, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

const EndpointVersion string = "/_version"

type httpVersionReponse struct {
	BuildVersion string `json:"build_version"`
}

func (s *Server) handlerVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet &&
			strings.EqualFold(r.URL.Path, EndpointVersion) {
			ver := httpVersionReponse{
				BuildVersion: s.buildVersion,
			}
			SendJSONPayload(r.Context(), w, "ok", ver)
			return
		}
		next.ServeHTTP(w, r)
	})
}

const EndpointHealthz string = "/_healthz"

func (s *Server) handlerReady(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet &&
			strings.EqualFold(r.URL.Path, EndpointHealthz) {

			if s.readyFunc == nil {
				sendText(r.Context(), w, 200, "undefined")
				return
			}

			if !s.readyFunc() {
				sendText(r.Context(), w, 500, "not ready")
				return
			}

			sendText(r.Context(), w, 200, ".")
			return
		}
		next.ServeHTTP(w, r)
	})
}

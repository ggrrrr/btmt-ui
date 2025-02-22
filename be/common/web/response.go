package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type AppResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Payload any    `json:"payload,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

func SendError(ctx context.Context, w http.ResponseWriter, e error) {
	var httpCode int = 500
	var msg string = ""
	var err error = e
	appError, ok := e.(*app.AppError)
	if ok {
		msg = appError.Msg()
		err = appError.Cause()
		httpCode = appError.Code()
	}
	sendJSON(ctx, w, httpCode, msg, err, nil)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	sendJSON(r.Context(), w, http.StatusNotFound, fmt.Sprintf("NotFound: [%s] %s", r.Method, r.URL.Path), nil, nil)
}

func methodNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	sendJSON(r.Context(), w, http.StatusMethodNotAllowed, fmt.Sprintf("NotFound: [%s] %s", r.Method, r.URL.Path), nil, nil)
}

// HTTP CODE 200
func SendJSONPayload(ctx context.Context, w http.ResponseWriter, msg string, payload any) {
	sendJSON(ctx, w, http.StatusOK, msg, nil, payload)
}

// HTTP CODE 500
func SendJSONSystemError(ctx context.Context, w http.ResponseWriter, msg string, err error, payload any) {
	sendJSON(ctx, w, http.StatusInternalServerError, msg, err, nil)
}

// HTTP CODE 400
func SendJSONErrorBadRequest(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	sendJSON(ctx, w, http.StatusBadRequest, msg, err, nil)
}

// HTTP CODE 401
func SendJSONUnauthorized(ctx context.Context, w http.ResponseWriter, msg string) {
	sendJSON(ctx, w, http.StatusUnauthorized, msg, nil, nil)
}

// HTTP CODE 403
func SendJSONForbidden(ctx context.Context, w http.ResponseWriter, msg string) {
	sendJSON(ctx, w, http.StatusForbidden, msg, nil, nil)
}

// // HTTP CODE 400
// func SendJSONErrorBadRequestWithBody(ctx context.Context, w http.ResponseWriter, msg string, err error, body any) {
// 	sendJSON(ctx, w, http.StatusBadRequest, msg, err, body)
// }

func sendJSON(ctx context.Context, w http.ResponseWriter, code int, msg string, err1 error, payload any) {
	traceID := ""
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		traceID = span.TraceID().String()
		w.Header().Add("X-Trace-ID", traceID)
	}

	errStr := ""
	if err1 != nil {
		errStr = err1.Error()
	}
	body := AppResponse{
		Code:    strconv.Itoa(code),
		Message: msg,
		Error:   errStr,
		Payload: payload,
		Trace:   traceID,
	}

	b, err := json.Marshal(body)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("unable to marshal response body")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("unable to write response")
	}
}

func sendText(ctx context.Context, w http.ResponseWriter, code int, text string) {
	w.Header().Add("Content-Type", "plain/text")
	w.WriteHeader(code)
	_, err := w.Write([]byte(text))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("unable to write response")
	}
}

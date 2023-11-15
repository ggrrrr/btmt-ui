package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"google.golang.org/grpc/codes"
)

type AppResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

func SendError(w http.ResponseWriter, e error) {
	var httpCode int = 500
	var msg string
	var err error

	appError, ok := e.(*app.AppError)
	if ok {
		msg = appError.Msg()
		err = appError.Err()
		switch appError.Code() {
		case codes.Internal:
			httpCode = 500
		case codes.InvalidArgument:
			httpCode = 400
		case codes.Unauthenticated:
			httpCode = 401
		case codes.PermissionDenied:
			httpCode = 403
		default:
			httpCode = 500
		}
	}
	send(w, httpCode, msg, err, nil)
}

func MethodNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	send(w, http.StatusNotFound, fmt.Sprintf("StatusNotFound: [%s] %s", r.Method, r.URL.Path), nil, nil)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	send(w, http.StatusMethodNotAllowed, fmt.Sprintf("StatusMethodNotAllowed: [%s] %s", r.Method, r.URL.Path), nil, nil)
}

func SendPayload(w http.ResponseWriter, msg string, payload any) {
	send(w, http.StatusOK, msg, nil, payload)
}

func SendSystemError(w http.ResponseWriter, msg string, err error, payload any) {
	send(w, http.StatusInternalServerError, msg, err, nil)
}

func SendErrorBadRequest(w http.ResponseWriter, msg string, err error) {
	send(w, http.StatusBadRequest, msg, err, nil)
}

func DecodeJsonRequest(r *http.Request, payload any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return app.ErrorSystem("unable to load http body", err)
	}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&payload)
	if err != nil {
		return app.ErrorBadRequest("bad json", err)
	}
	return nil
}

func send(w http.ResponseWriter, code int, msg string, err1 error, payload any) {
	errStr := ""
	if err1 != nil {
		errStr = err1.Error()
	}
	body := AppResponse{
		Code:    fmt.Sprintf("%d", code),
		Message: msg,
		Error:   errStr,
		Payload: payload,
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Printf("unable to write reponse body(%v) error: %v", body, err)
	}
	w.WriteHeader(code)
	w.Write(b)
}

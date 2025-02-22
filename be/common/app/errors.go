package app

import (
	"fmt"
	"net/http"
	// grpcCodes "google.golang.org/grpc/codes"
	// grpcStatus "google.golang.org/grpc/status"
)

var (
	ErrAuthUnauthenticated = UnauthenticatedError("please login", nil)
	ErrForbidden           = PermissionDeniedError("forbidden", nil)
	ErrTeapot              = SystemError("teapot", nil)
)

type AppError struct {
	msg  string
	err  error
	code int
	// grpcCode grpcCodes.Code
}

func (e *AppError) Cause() error {
	return e.err
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Error() string {
	errStr := ""
	if e.err != nil {
		return fmt.Sprintf("[%d]: %s -> %v", e.code, e.msg, errStr)
	}
	return fmt.Sprintf("[%d]: %s", e.code, e.msg)
}

// HTTP CODE 500
func SystemError(msg string, err error) error {
	return &AppError{
		code: http.StatusInternalServerError,
		msg:  msg,
		err:  err,
	}
}

// HTTP CODE 400
func BadRequestError(msg string, err error) error {
	return &AppError{
		code: http.StatusBadRequest,
		msg:  msg,
		err:  err,
	}
}

// HTTP code 401
func UnauthenticatedError(msg string, err error) error {
	return &AppError{
		code: http.StatusUnauthorized,
		msg:  msg,
		err:  err,
	}
}

// HTTP code 403
func PermissionDeniedError(msg string, err error) error {
	return &AppError{
		code: http.StatusForbidden,
		msg:  msg,
		err:  err,
	}
}

// HTTP code 404
func ItemNotFoundError(itemName, itemID string) error {
	return &AppError{
		code: http.StatusNotFound,
		msg:  fmt.Sprintf("%s with ID:%s not found", itemName, itemID),
		err:  nil,
	}
}

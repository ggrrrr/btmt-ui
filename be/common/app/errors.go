package app

import (
	"fmt"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

var (
	ErrAuthUnauthenticated = UnauthenticatedError("please login", nil)
	ErrForbidden           = PermissionDeniedError("forbidden", nil)
	ErrTeapot              = SystemError("teapot", nil)
)

type AppError struct {
	msg      string
	err      error
	grpcCode grpcCodes.Code
}

func (e *AppError) Cause() error {
	return e.err
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) Code() grpcCodes.Code {
	return e.grpcCode
}

func (e *AppError) Error() string {
	errStr := ""
	if e.err != nil {
		return fmt.Sprintf("[%d]: %s -> %v", e.grpcCode, e.msg, errStr)
	}
	return fmt.Sprintf("[%d]: %s", e.grpcCode, e.msg)
}

func ToGrpcError(e error) error {
	if e == nil {
		return nil
	}
	appError, ok := e.(*AppError)
	if ok {
		return grpcStatus.Error(appError.grpcCode, appError.Msg())
	}
	return grpcStatus.Error(grpcCodes.Internal, e.Error())
}

// HTTP CODE 500
func SystemError(msg string, err error) error {
	return &AppError{
		grpcCode: grpcCodes.Internal,
		msg:      msg,
		err:      err,
	}
}

// HTTP CODE 400
func BadRequestError(msg string, err error) error {
	return &AppError{
		grpcCode: grpcCodes.InvalidArgument,
		msg:      msg,
		err:      err,
	}
}

// HTTP code 401
func UnauthenticatedError(msg string, err error) error {
	return &AppError{
		grpcCode: grpcCodes.Unauthenticated,
		msg:      msg,
		err:      err,
	}
}

// HTTP code 403
func PermissionDeniedError(msg string, err error) error {
	return &AppError{
		grpcCode: grpcCodes.PermissionDenied,
		msg:      msg,
		err:      err,
	}
}

// HTTP code 404
func ItemNotFoundError(itemName, itemID string) error {
	return &AppError{
		grpcCode: grpcCodes.NotFound,
		msg:      fmt.Sprintf("%s with ID:%s not found", itemName, itemID),
		err:      nil,
	}
}

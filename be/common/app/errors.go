package app

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthUnauthenticated = ErrorUnauthenticated("please login", nil)
	ErrForbidden           = ErrorDenied("forbidden", nil)
	ErrTeapot              = ErrorSystem("teapot", nil)
)

type AppError struct {
	msg      string
	err      error
	grpcCode codes.Code
}

func (e *AppError) Err() error {
	return e.err
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) Code() codes.Code {
	return e.grpcCode
}

func (e *AppError) Error() string {
	errStr := ""
	if e.err != nil {
		errStr = e.err.Error()
	}
	return fmt.Sprintf("[%d]: %s %s", e.grpcCode, e.msg, errStr)
}

func ToGrpcError(e error) error {
	if e == nil {
		return nil
	}
	appError, ok := e.(*AppError)
	if ok {
		return status.Error(appError.grpcCode, appError.Msg())
	}
	return status.Error(codes.Internal, e.Error())
}

func ErrorSystem(msg string, err error) error {
	return &AppError{
		grpcCode: codes.Internal,
		msg:      msg,
		err:      err,
	}
}

func ErrorBadRequest(msg string, err error) error {
	return &AppError{
		grpcCode: codes.InvalidArgument,
		msg:      msg,
		err:      err,
	}
}

func ErrorUnauthenticated(msg string, err error) error {
	return &AppError{
		grpcCode: codes.Unauthenticated,
		msg:      msg,
		err:      err,
	}
}

func ErrorDenied(msg string, err error) error {
	return &AppError{
		grpcCode: codes.PermissionDenied,
		msg:      msg,
		err:      err,
	}
}

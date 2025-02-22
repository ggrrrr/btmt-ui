package app

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

var (
	ErrAuthEmailEmpty    = app.BadRequestError("username empty", nil)
	ErrAuthPasswdEmpty   = app.BadRequestError("password empty", nil)
	ErrAuthEmailNotFound = app.BadRequestError("account not found", nil)
	ErrAuthUserPassword  = app.UnauthenticatedError("wrong username/password", nil)
	ErrAuthEmailLocked   = app.UnauthenticatedError("account locked", nil)

	ErrAuthMultipleEmail = app.SystemError("multiple emails", nil)
)

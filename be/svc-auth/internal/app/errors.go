package app

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

var (
	ErrAuthEmailEmpty    = app.BadRequestError("email empty", nil)
	ErrAuthPasswdEmpty   = app.BadRequestError("password empty", nil)
	ErrAuthEmailNotFound = app.BadRequestError("email not found", nil)
	ErrAuthBadPassword   = app.UnauthenticatedError("wrong password", nil)
	ErrAuthEmailLocked   = app.UnauthenticatedError("emails locked", nil)

	ErrAuthMultipleEmail = app.SystemError("multiple emails", nil)
)

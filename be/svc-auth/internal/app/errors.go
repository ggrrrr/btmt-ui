package app

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

var (
	errAuthEmailEmpty    = app.BadRequestError("email empty", nil)
	errAuthPasswdEmpty   = app.BadRequestError("password empty", nil)
	errAuthEmailNotFound = app.BadRequestError("email not found", nil)
	errAuthBadPassword   = app.UnauthenticatedError("wrong password", nil)
	errAuthEmailLocked   = app.UnauthenticatedError("emails locked", nil)

	errAuthMultipleEmail = app.SystemError("multiple emails", nil)
)

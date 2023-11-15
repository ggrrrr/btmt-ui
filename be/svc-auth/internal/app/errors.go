package app

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

var (
	ErrAuthEmailEmpty    = app.ErrorBadRequest("email empty", nil)
	ErrAuthPasswdEmpty   = app.ErrorBadRequest("password empty", nil)
	ErrAuthEmailNotFound = app.ErrorUnauthenticated("email not found", nil)
	ErrAuthBadPassword   = app.ErrorUnauthenticated("wrong password", nil)
	ErrAuthEmailLocked   = app.ErrorUnauthenticated("emails locked", nil)

	ErrAuthMultipleEmail = app.ErrorSystem("multiple emails", nil)
)

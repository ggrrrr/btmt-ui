package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func (ap *application) LoginPasswd(ctx context.Context, email, passwd string) (app.Result[AuthToken], error) {
	if email == "" {
		return app.Result[AuthToken]{}, ErrAuthEmailEmpty
	}
	if passwd == "" {
		return app.Result[AuthToken]{}, ErrAuthPasswdEmpty
	}
	auth, err := ap.findEmail(ctx, email)
	if err != nil {
		logger.Error(err).Msg("ap.findEmail")
		return app.Result[AuthToken]{}, app.SystemError("failed to fetch email", err)
	}

	if auth == nil {
		return app.Result[AuthToken]{}, ErrAuthEmailNotFound
	}

	if !canLogin(auth) {
		return app.Result[AuthToken]{}, ErrAuthEmailLocked
	}

	if !checkPasswordHash(passwd, string(auth.Passwd)) {
		return app.Result[AuthToken]{}, ErrAuthBadPassword
	}

	jwt, err := ap.signer.Sign(auth.ToAuthInfo(roles.SystemTenant))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.signer.Sign")
		return app.Result[AuthToken]{}, app.SystemError("Unable to sign, please try again later", nil)
	}
	return app.ResultWithPayload[AuthToken]("ok", AuthToken(jwt)), nil
}

func (ap *application) TokenValidate(ctx context.Context) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if authInfo.User == "" {
		return app.ErrAuthUnauthenticated
	}
	auth, err := ap.findEmail(ctx, authInfo.User)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Validate")
		return app.SystemError("failed to fetch email", err)
	}
	if auth == nil {
		return ErrAuthEmailNotFound
	}
	if !canLogin(auth) {
		return ErrAuthEmailLocked
	}
	return nil
}

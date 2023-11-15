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
		logger.Log().Debug().Err(err).Msg("failed to fetch email")
		return app.Result[AuthToken]{}, app.ErrorSystem("failed to fetch email", err)
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

	jwt, err := ap.signer.Sign(auth.ToAuthInfo())
	if err != nil {
		return app.Result[AuthToken]{}, err
	}
	return app.ResultPayload[AuthToken]("ok", AuthToken(jwt)), nil
}

func (ap *application) Validate(ctx context.Context) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if authInfo.User == "" {
		return app.ErrAuthUnauthenticated
	}
	return nil
}

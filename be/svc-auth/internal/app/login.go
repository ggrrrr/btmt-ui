package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (ap *Application) LoginPasswd(ctx context.Context, email, passwd string) (ddd.AuthToken, error) {
	var err error
	ctx, span := logger.Span(ctx, "LoginPasswd", nil)
	defer func() {
		span.End(err)
	}()

	if email == "" {
		err = ErrAuthEmailEmpty
		return ddd.AuthToken{}, err
	}
	if passwd == "" {
		err = ErrAuthPasswdEmpty
		return ddd.AuthToken{}, err
	}

	auth, err := ap.findEmail(ctx, email)
	if err != nil {
		logger.Error(err).Msg("ap.findEmail")
		return ddd.AuthToken{}, err
	}

	if auth == nil {
		err = ErrAuthEmailNotFound
		return ddd.AuthToken{}, err
	}

	if !canLogin(auth) {
		err = ErrAuthEmailLocked
		return ddd.AuthToken{}, err
	}

	if !checkPasswordHash(passwd, string(auth.Passwd)) {
		err = ErrAuthBadPassword
		return ddd.AuthToken{}, err
	}

	jwt, expiresAt, err := ap.signer.Sign(auth.ToAuthInfo(roles.SystemRealm))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.signer.Sign")
		return ddd.AuthToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	return ddd.AuthToken{
		Token:     jwt,
		ExpiresAt: expiresAt,
	}, nil
}

func (ap *Application) TokenValidate(ctx context.Context) (err error) {
	ctx, span := logger.Span(ctx, "TokenValidate", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if authInfo.User == "" {
		err = app.ErrAuthUnauthenticated
		return
	}
	auth, err := ap.findEmail(ctx, authInfo.User)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Validate")
		return app.SystemError("failed to fetch email", err)
	}
	if auth == nil {
		err = ErrAuthEmailNotFound
		return
	}
	if !canLogin(auth) {
		err = ErrAuthEmailLocked
		return
	}
	return nil
}

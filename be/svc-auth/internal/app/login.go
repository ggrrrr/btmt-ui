package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (ap *Application) LoginPasswd(ctx context.Context, email, passwd string) (ddd.LoginToken, error) {
	var err error
	ctx, span := logger.Span(ctx, "LoginPasswd", nil)
	defer func() {
		span.End(err)
	}()

	if email == "" {
		err = ErrAuthEmailEmpty
		return ddd.LoginToken{}, err
	}
	if passwd == "" {
		err = ErrAuthPasswdEmpty
		return ddd.LoginToken{}, err
	}

	authPasswd, err := ap.findEmail(ctx, email)
	if err != nil {
		logger.Error(err).Msg("ap.findEmail")
		return ddd.LoginToken{}, err
	}

	if authPasswd == nil {
		err = ErrAuthEmailNotFound
		return ddd.LoginToken{}, err
	}

	if !canLogin(authPasswd) {
		err = ErrAuthEmailLocked
		return ddd.LoginToken{}, err
	}

	if !checkPasswordHash(passwd, string(authPasswd.Passwd)) {
		err = ErrAuthBadPassword
		return ddd.LoginToken{}, err
	}

	currentAuthInfo := roles.AuthInfoFromCtx(ctx)

	authInfo := authPasswd.ToAuthInfo(currentAuthInfo.Device, roles.SystemRealm)

	accessToken, expiresAt, err := ap.signer.Sign(ctx, ap.accessTokenTTL, authInfo)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	refreshRole := roles.AuthInfo{
		User:  authPasswd.Email,
		Realm: roles.SystemRealm,
		Roles: []string{authpb.AuthSvc_TokenRefresh_FullMethodName},
		ID:    authInfo.ID,
	}

	// asd := authpb.AuthSvc_To/
	refreshToken, refreshExpiresAt, err := ap.signer.Sign(ctx, ap.refreshTokenTTL, refreshRole)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	err = ap.historyRepo.SaveHistory(ctx, authInfo, authpb.AuthSvc_LoginPasswd_FullMethodName)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.authRepo.SaveHistory")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	return ddd.LoginToken{
		ID:    authInfo.ID,
		Email: email,
		AccessToken: ddd.AuthToken{
			Value:     accessToken,
			ExpiresAt: expiresAt,
		},
		RefreshToken: ddd.AuthToken{
			Value:     refreshToken,
			ExpiresAt: refreshExpiresAt,
		},
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

func (ap *Application) TokenRefresh(ctx context.Context) (loginToken ddd.LoginToken, err error) {
	ctx, span := logger.Span(ctx, "TokenRefresh", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)

	err = ap.appPolices.CanDo(authInfo.Realm, authpb.AuthSvc_TokenRefresh_FullMethodName, authInfo)
	if err != nil {
		return loginToken, app.PermissionDeniedError("token refresh roles", err)
	}

	auth, err := ap.findEmail(ctx, authInfo.User)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("TokenRefresh")
		return loginToken, app.SystemError("failed to fetch email", err)
	}
	if auth == nil {
		err = ErrAuthEmailNotFound
		return
	}
	if !canLogin(auth) {
		err = ErrAuthEmailLocked
		return
	}

	history, err := ap.historyRepo.GetHistory(ctx, authInfo.ID)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("TokenRefresh")
		return loginToken, app.SystemError("failed to fetch login history", err)
	}

	if history == nil {
		return loginToken, app.UnauthenticatedError("please login", nil)

	}

	newAuthInfo := auth.ToAuthInfo(authInfo.Device, authInfo.Realm)

	jwtValue, expiresAt, err := ap.signer.Sign(ctx, ap.accessTokenTTL, newAuthInfo)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	loginToken = ddd.LoginToken{
		ID:    authInfo.ID,
		Email: authInfo.User,
		AccessToken: ddd.AuthToken{
			Value:     jwtValue,
			ExpiresAt: expiresAt,
		},
	}
	return
}

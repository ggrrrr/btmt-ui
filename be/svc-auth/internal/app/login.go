package app

import (
	"context"
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	authpb "github.com/ggrrrr/btmt-ui/be/svc-auth/authpb/v1"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (ap *Application) LoginPasswd(ctx context.Context, username, passwd string) (ddd.LoginToken, error) {
	var err error
	ctx, span := ap.otelTracer.SpanWithAttributes(ctx, "LoginPasswd", slog.String("username", username))
	defer func() {
		span.End(err)
	}()

	if username == "" {
		err = ErrAuthEmailEmpty
		return ddd.LoginToken{}, err
	}
	if passwd == "" {
		err = ErrAuthPasswdEmpty
		return ddd.LoginToken{}, err
	}

	authPasswd, err := ap.findEmail(ctx, username)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "ap.findEmail")
		return ddd.LoginToken{}, ErrAuthUserPassword
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
		err = ErrAuthUserPassword
		return ddd.LoginToken{}, err
	}

	currentAuthInfo := roles.AuthInfoFromCtx(ctx)

	authInfo := authPasswd.ToAuthInfo(currentAuthInfo.Device, roles.SystemRealm)

	accessToken, expiresAt, err := ap.signer.Sign(ctx, ap.accessTokenTTL, authInfo)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	refreshRole := roles.AuthInfo{
		Subject: authPasswd.Subject,
		// TODO sudo like...
		// AdminSubject: "",
		Realm: roles.SystemRealm,
		Roles: []string{authpb.AuthSvc_TokenRefresh_FullMethodName},
		ID:    authInfo.ID,
	}

	refreshToken, refreshExpiresAt, err := ap.signer.Sign(ctx, ap.refreshTokenTTL, refreshRole)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	err = ap.historyRepo.SaveHistory(ctx, authInfo, authpb.AuthSvc_LoginPasswd_FullMethodName)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "ap.authRepo.SaveHistory")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	return ddd.LoginToken{
		ID:           authInfo.ID,
		Subject:      username,
		AdminSubject: "",
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
	authInfo := roles.AuthInfoFromCtx(ctx)
	ctx, span := ap.otelTracer.SpanWithAttributes(ctx, "TokenValidate")
	defer span.End(err)

	if authInfo.Subject == "" {
		err = app.ErrAuthUnauthenticated
		return
	}
	auth, err := ap.findEmail(ctx, authInfo.Subject)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Validate")
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
	authInfo := roles.AuthInfoFromCtx(ctx)
	ctx, span := ap.otelTracer.SpanWithAttributes(ctx, "TokenRefresh")
	defer span.End(err)

	err = ap.appPolices.CanDo(authInfo.Realm, authpb.AuthSvc_TokenRefresh_FullMethodName, authInfo)
	if err != nil {
		return loginToken, app.PermissionDeniedError("token refresh roles", err)
	}

	auth, err := ap.findEmail(ctx, authInfo.Subject)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "TokenRefresh")
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
		log.Log().ErrorCtx(ctx, err, "TokenRefresh")
		return loginToken, app.SystemError("failed to fetch login history", err)
	}

	if history == nil {
		return loginToken, app.UnauthenticatedError("please login", nil)

	}

	newAuthInfo := auth.ToAuthInfo(authInfo.Device, authInfo.Realm)

	jwtValue, expiresAt, err := ap.signer.Sign(ctx, ap.accessTokenTTL, newAuthInfo)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "ap.signer.Sign")
		return ddd.LoginToken{}, app.SystemError("Unable to sign, please try again later", nil)
	}

	loginToken = ddd.LoginToken{
		ID:      authInfo.ID,
		Subject: authInfo.Subject,
		AccessToken: ddd.AuthToken{
			Value:     jwtValue,
			ExpiresAt: expiresAt,
		},
	}
	return
}

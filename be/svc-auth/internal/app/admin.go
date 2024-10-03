package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (a *Application) UserCreate(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.Span(ctx, "UserCreate", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserCreate_FullMethodName, authInfo)
	if err != nil {
		return err
	}

	if auth.Email == "" {
		err = ErrAuthEmailEmpty
		return
	}

	if auth.Passwd == "" {
		err = ErrAuthPasswdEmpty
		return err
	}

	logger.InfoCtx(ctx).Any("email", auth.Email).Msg("CreateAuth")
	if auth.Passwd != "" {
		cryptPasswd, err := HashPassword(string(auth.Passwd))
		if err != nil {
			logger.ErrorCtx(ctx, err).Msg("UserCreate.HashPassword")
			return app.SystemError("Create password", err)
		}
		auth.Passwd = cryptPasswd
	}
	err = a.authRepo.Save(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

func (ap *Application) Get(ctx context.Context, email string) (result app.Result[ddd.AuthPasswd], err error) {
	ctx, span := logger.Span(ctx, "Get", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err = ap.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserList_FullMethodName, authInfo); err != nil {
		return
	}

	auth, err := ap.findEmail(ctx, email)
	if err != nil {
		logger.Error(err).Msg("ap.findEmail")
		return app.Result[ddd.AuthPasswd]{}, err
	}

	if auth == nil {
		err = ErrAuthEmailNotFound
		return app.Result[ddd.AuthPasswd]{}, err
	}

	return app.ResultWithPayload[ddd.AuthPasswd]("ok", *auth), nil

}

func (ap *Application) UserList(ctx context.Context) (result app.Result[[]ddd.AuthPasswd], err error) {
	ctx, span := logger.Span(ctx, "UserList", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err = ap.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserList_FullMethodName, authInfo); err != nil {
		return
	}

	out, err := ap.authRepo.List(ctx, nil)
	if err != nil {
		return
	}

	logger.InfoCtx(ctx).Msg("UserList")
	logger.DebugCtx(ctx).
		Any("list", out).
		Msg("UserList")
	return app.ResultWithPayload[[]ddd.AuthPasswd]("ok", out), nil
}

func (ap *Application) UserUpdate(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.Span(ctx, "UserUpdate", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = ap.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserUpdate_FullMethodName, authInfo)
	if err != nil {
		return
	}

	list, err := ap.authRepo.Get(ctx, auth.Email)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("authRepo.Get")
		return
	}
	if len(list) == 0 {
		err = ErrAuthEmailNotFound
		return
	}
	update := list[0]
	update.Status = auth.Status
	update.SystemRoles = auth.SystemRoles
	return ap.authRepo.Update(ctx, update)
}

func (a *Application) UserChangePasswd(ctx context.Context, email, oldPasswd, newPasswd string) (err error) {
	ctx, span := logger.Span(ctx, "UserChangePasswd", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserChangePasswd_FullMethodName, authInfo)
	if err != nil {
		return
	}

	rec, err := a.findEmail(ctx, email)
	if err != nil {
		return
	}

	if rec == nil {
		err = ErrAuthEmailNotFound
		return
	}

	if !checkPasswordHash(oldPasswd, string(rec.Passwd)) {
		err = ErrAuthBadPassword
		return
	}

	var cryptPasswd string
	if newPasswd != "" {
		cryptPasswd, err = HashPassword(newPasswd)
		if err != nil {
			return app.SystemError("HashPassword", err)
		}
	}
	logger.InfoCtx(ctx).Any("email", email).Msg("UpdatePasswd")
	err = a.authRepo.UpdatePassword(ctx, email, cryptPasswd)
	if err != nil {
		return
	}
	return err
}

package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (a *application) UserCreate(ctx context.Context, auth ddd.AuthPasswd) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserCreate_FullMethodName, authInfo); err != nil {
		return err
	}
	if auth.Email == "" {
		return app.BadRequestError("email empty", nil)
	}
	if auth.Passwd == "" {
		return app.BadRequestError("password empty", nil)
	}
	logger.InfoCtx(ctx).Any("email", auth.Email).Msg("CreateAuth")
	if auth.Passwd != "" {
		cryptPasswd, err := HashPassword(string(auth.Passwd))
		if err != nil {
			return err
		}
		auth.Passwd = cryptPasswd
	}
	err := a.authRepo.Save(ctx, auth)
	return err
}

func (ap *application) UserList(ctx context.Context) (app.Result[[]ddd.AuthPasswd], error) {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := ap.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserList_FullMethodName, authInfo); err != nil {
		return app.Result[[]ddd.AuthPasswd]{}, err
	}
	out, err := ap.authRepo.List(ctx)
	if err != nil {
		return app.Result[[]ddd.AuthPasswd]{}, err
	}
	logger.InfoCtx(ctx).Msg("ListAuth")
	return app.ResultWithPayload[[]ddd.AuthPasswd]("ok", out), err
}

func (ap *application) UserUpdate(ctx context.Context, auth ddd.AuthPasswd) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := ap.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserUpdate_FullMethodName, authInfo); err != nil {
		return err
	}
	list, err := ap.authRepo.Get(ctx, auth.Email)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return app.BadRequestError("email not found", nil)
	}
	update := list[0]
	update.Status = auth.Status
	update.SystemRoles = auth.SystemRoles
	return ap.authRepo.Update(ctx, update)
}

func (a *application) UserChangePasswd(ctx context.Context, email, oldPasswd, newPasswd string) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(roles.SystemTenant, authpb.AuthSvc_UserChangePasswd_FullMethodName, authInfo); err != nil {
		return err
	}
	rec, err := a.findEmail(ctx, email)
	if err != nil {
		return err
	}
	if rec == nil {
		return ErrAuthEmailNotFound
	}
	if !checkPasswordHash(oldPasswd, string(rec.Passwd)) {
		return ErrAuthBadPassword
	}

	var cryptPasswd string
	if newPasswd != "" {
		cryptPasswd, err = HashPassword(newPasswd)
		if err != nil {
			return err
		}
	}
	logger.InfoCtx(ctx).Any("email", email).Msg("UpdatePasswd")
	err = a.authRepo.UpdatePassword(ctx, email, cryptPasswd)
	if err != nil {
		return err
	}
	return err
}

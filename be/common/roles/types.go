package roles

import (
	"context"

	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

const AuthSchemeBearer string = "Bearer"

type (
	authInfoCtxKeyType struct{}

	AuthInfo struct {
		// User name or system name
		Subject string
		// For `sudo` like behavior
		ForSubject string
		// Domain name
		Realm string
		// List of domain roles
		Roles []string
		// List of system roles
		SystemRoles []string
		// Info of the agent device ( browser,app,service,etc...)
		Device app.Device
		// Unique ID of the token
		ID uuid.UUID
	}
)

func CreateSystemAdminUser(tenant string, subject string, device app.Device) AuthInfo {
	return AuthInfo{
		Realm:       tenant,
		Subject:     subject,
		Roles:       []string{RoleAdmin},
		SystemRoles: []string{RoleAdmin},
		Device:      device,
	}
}

func CtxWithAuthInfo(ctx context.Context, authInfo AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoCtxKeyType{}, authInfo)
}

func AuthInfoFromCtx(ctx context.Context) AuthInfo {
	value, ok := ctx.Value(authInfoCtxKeyType{}).(AuthInfo)
	if !ok {
		return AuthInfo{}
	}
	return value
}

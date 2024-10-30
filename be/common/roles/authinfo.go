package roles

import (
	"context"

	"github.com/google/uuid"
)

type (
	AuthInfo struct {
		Subject     string
		Realm       string
		Roles       []string
		SystemRoles []string
		Device      Device
		ID          uuid.UUID
	}
)

func CreateSystemAdminUser(tenant string, subject string, device Device) AuthInfo {
	return AuthInfo{
		Realm:       tenant,
		Subject:     subject,
		Roles:       []string{RoleAdmin},
		SystemRoles: []string{RoleAdmin},
		Device:      device,
	}
}

func CtxWithAuthInfo(ctx context.Context, authInfo AuthInfo) context.Context {
	return context.WithValue(ctx, ctxKeyType{}, authInfo)
}

func AuthInfoFromCtx(ctx context.Context) AuthInfo {
	value, ok := ctx.Value(ctxKeyType{}).(AuthInfo)
	if !ok {
		return AuthInfo{}
	}
	return value
}

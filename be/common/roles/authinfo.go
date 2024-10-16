package roles

import "context"

type (
	AuthInfo struct {
		User        string
		Realm       string
		Roles       []string
		SystemRoles []string
		Device      Device
	}
)

func CreateSystemAdminUser(tenant string, user string, device Device) AuthInfo {
	return AuthInfo{
		Realm:       tenant,
		User:        user,
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

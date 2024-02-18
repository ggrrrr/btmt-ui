package roles

import "context"

type (
	AuthInfo struct {
		User        string
		Tenant      Tenant
		Roles       []RoleName
		SystemRoles []RoleName
		Device      Device
	}
)

func CreateSystemAdminUser(tenant Tenant, user string, device Device) AuthInfo {
	return AuthInfo{
		Tenant:      tenant,
		User:        user,
		Roles:       []RoleName{RoleAdmin},
		SystemRoles: []RoleName{RoleAdmin},
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

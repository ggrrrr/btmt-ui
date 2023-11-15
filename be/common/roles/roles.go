package roles

import (
	"context"
)

const AuthSchemeBeaerer string = "Bearer"

type (
	ctxKeyType struct{}

	RoleName string

	AuthCredentials string

	Authorization struct {
		AuthScheme      string          // Basic,Beaerer,
		AuthCredentials AuthCredentials // JWT TOKEN OR OTHER secret data
	}
	UserRequest struct {
		FullMethod    string
		Device        Device
		Authorization Authorization
	}

	Device struct {
		RemoteAddr string
		DeviceInfo string
	}

	AuthInfo struct {
		User   string
		Roles  []RoleName
		Device Device
	}
)

const (
	RoleAdmin RoleName = "admin"
)

func CreateAdminUser(user string, device Device) AuthInfo {
	return AuthInfo{
		User:   user,
		Roles:  []RoleName{RoleAdmin},
		Device: device,
	}
}

func HasRole(role RoleName, roles []RoleName) bool {
	for r := range roles {
		if roles[r] == role {
			return true
		}
	}
	return false
}

func isSystemAdmin(authInfo AuthInfo) bool {
	return HasRole(RoleAdmin, authInfo.Roles)
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

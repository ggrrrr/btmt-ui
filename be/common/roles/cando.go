package roles

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

type (
	canDo struct {
	}

	AppPolices interface {
		CanDo(tenant string, FullMethodName string, authInfo AuthInfo) error
	}
)

const (
	RoleAdmin        string = "admin"
	RoleTokenRefresh string = "token.refresh"
	SystemRealm      string = "localhost"
)

var _ (AppPolices) = (*canDo)(nil)

func NewAppPolices() *canDo {
	return &canDo{}
}

func (*canDo) CanDo(tenant string, fullMethodName string, authInfo AuthInfo) error {
	if authInfo.Subject == "" {
		return app.ErrAuthUnauthenticated
	}
	if isSystemAdmin(authInfo) {
		return nil
	}

	if isAdmin(authInfo) {
		return nil
	}
	if tenant == authInfo.Realm {
		if hasRole(fullMethodName, authInfo.Roles) {
			return nil
		}
	}
	return app.ErrForbidden
}

func isAdmin(authInfo AuthInfo) bool {
	return hasRole(RoleAdmin, authInfo.Roles)
}

func isSystemAdmin(authInfo AuthInfo) bool {
	return hasRole(RoleAdmin, authInfo.SystemRoles)
}

func hasRole(role string, roles []string) bool {
	for r := range roles {
		if roles[r] == role {
			return true
		}
	}
	return false
}

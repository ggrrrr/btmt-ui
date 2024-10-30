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
		if HasRole(fullMethodName, authInfo.Roles) {
			return nil
		}
	}
	return app.ErrForbidden
}

func isAdmin(authInfo AuthInfo) bool {
	return HasRole(RoleAdmin, authInfo.Roles)
}

func isSystemAdmin(authInfo AuthInfo) bool {
	return HasRole(RoleAdmin, authInfo.SystemRoles)
}

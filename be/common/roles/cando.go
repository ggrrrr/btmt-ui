package roles

import (
	"github.com/ggrrrr/btmt-ui/be/common/app"
)

type (
	canDo struct {
	}

	AppPolices interface {
		CanDo(tenant Tenant, FullMethodName string, authInfo AuthInfo) error
	}
)

var _ (AppPolices) = (*canDo)(nil)

func NewAppPolices() *canDo {
	return &canDo{}
}

func (*canDo) CanDo(tenant Tenant, FullMethodName string, authInfo AuthInfo) error {
	if authInfo.User == "" {
		return app.ErrAuthUnauthenticated
	}
	if tenant == SystemTenant {
		if isSystemAdmin(authInfo) {
			return nil
		}
		return app.ErrForbidden
	}
	if tenant != authInfo.Tenant {
		return app.ErrForbidden
	}

	if isAdmin(authInfo) {
		return nil
	}
	return app.ErrForbidden
}

func isAdmin(authInfo AuthInfo) bool {
	return HasRole(RoleAdmin, authInfo.Roles)
}

func isSystemAdmin(authInfo AuthInfo) bool {
	return HasRole(RoleAdmin, authInfo.SystemRoles)
}

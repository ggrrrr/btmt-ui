package roles

import "github.com/ggrrrr/btmt-ui/be/common/app"

type (
	canDo struct {
	}

	AppPolices interface {
		CanDo(FullMethodName string, authInfo AuthInfo) error
	}
)

var _ (AppPolices) = (*canDo)(nil)

func NewAppPolices() *canDo {
	return &canDo{}
}

func (*canDo) CanDo(FullMethodName string, authInfo AuthInfo) error {
	if isSystemAdmin(authInfo) {
		return nil
	}
	if authInfo.User == "" {
		return app.ErrAuthUnauthenticated
	}
	return app.ErrForbidden
}

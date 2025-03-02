package token

import (
	"log/slog"
	"strings"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	verifier_mock struct {
	}
)

var _ (Verifier) = (*verifier_mock)(nil)

func NewVerifierMock() *verifier_mock {
	log.Log().Warn(nil, "NewVerifierMock")
	return &verifier_mock{}
}

func (*verifier_mock) Verify(auth app.AuthData) (roles.AuthInfo, error) {
	log.Log().Warn(nil, "NewVerifierMock", slog.Any("", auth))
	if auth.AuthScheme != roles.AuthSchemeBearer {
		log.Log().Error(ErrJwtBadScheme, "NewVerifierMock", slog.Any("", auth))
		return roles.AuthInfo{}, ErrJwtBadScheme
	}

	return roles.AuthInfo{
		Subject:     strings.Split(string(auth.AuthToken), " ")[0],
		Roles:       []string{"admin"},
		SystemRoles: []string{"admin"},
		Realm:       roles.SystemRealm,
	}, nil
}

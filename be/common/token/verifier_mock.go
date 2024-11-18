package token

import (
	"fmt"
	"strings"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	verifier_mock struct {
	}
)

var _ (Verifier) = (*verifier_mock)(nil)

func NewVerifierMock() *verifier_mock {
	logger.Warn().Msg("NewVerifierMock")
	return &verifier_mock{}
}

func (*verifier_mock) Verify(auth app.AuthData) (roles.AuthInfo, error) {
	logger.Warn().Any("token", auth).Msg("NewVerifierMock.Verify")
	if auth.AuthScheme != "mock" {
		logger.Error(fmt.Errorf("AuthScheme is not mock")).Any("token", auth).Str("AuthScheme", auth.AuthScheme).Msg("NewVerifierMock.Verify")
		return roles.AuthInfo{}, ErrJwtBadScheme
	}

	return roles.AuthInfo{
		Subject:     strings.Split(string(auth.AuthToken), " ")[0],
		Roles:       []string{"admin"},
		SystemRoles: []string{"admin"},
		Realm:       roles.SystemRealm,
	}, nil
}

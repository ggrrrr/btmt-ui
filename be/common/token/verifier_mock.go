package token

import (
	"fmt"
	"strings"

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

func (*verifier_mock) Verify(auth roles.Authorization) (roles.AuthInfo, error) {
	logger.Warn().Any("token", auth).Msg("NewVerifierMock.Verify")
	if auth.AuthScheme != "mock" {
		logger.Error(fmt.Errorf("AuthScheme is not mock")).Any("token", auth).Str("AuthScheme", auth.AuthScheme).Msg("NewVerifierMock.Verify")
		return roles.AuthInfo{}, ErrJwtBadScheme
	}

	return roles.AuthInfo{
		User:        strings.Split(string(auth.AuthCredentials), " ")[0],
		Roles:       []roles.RoleName{"admin"},
		SystemRoles: []roles.RoleName{"admin"},
		Tenant:      roles.SystemTenant,
	}, nil
}

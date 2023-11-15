package token

import (
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
	logger.Log().Warn().Msg("NewVerifierMock")
	return &verifier_mock{}
}

func (*verifier_mock) Verify(auth roles.Authorization) (roles.AuthInfo, error) {
	logger.Log().Warn().Any("token", auth).Msg("NewVerifierMock.Verify")
	if auth.AuthScheme != "mock" {
		logger.Log().Error().Any("token", auth).Str("AuthScheme", auth.AuthScheme).Msg("NewVerifierMock.Verify")
		return roles.AuthInfo{}, ErrJwtBadScheme
	}

	return roles.AuthInfo{
		User:  strings.Split(string(auth.AuthCredentials), " ")[0],
		Roles: []roles.RoleName{"admin"},
	}, nil
}

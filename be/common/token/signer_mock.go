package token

import (
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	mocker struct{}
)

var _ (Signer) = (*mocker)(nil)

func NewSignerMock() *mocker {
	logger.Warn().Msg("NewSignerMock")
	return &mocker{}
}

func (*mocker) Sign(claims roles.AuthInfo) (string, error) {
	return "ok", nil
}

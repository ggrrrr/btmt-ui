package token

import (
	"fmt"
	"time"

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

func (*mocker) Sign(claims roles.AuthInfo) (string, time.Time, error) {
	if claims.User == "" {
		return "", time.Time{}, fmt.Errorf("mock user is empty")
	}
	return "ok", time.Now().UTC().Add(1 * time.Hour), nil
}

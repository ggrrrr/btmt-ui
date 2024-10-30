package token

import (
	"context"
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

func (*mocker) Sign(_ context.Context, ttl time.Duration, claims roles.AuthInfo) (string, time.Time, error) {
	if claims.Subject == "" {
		return "", time.Time{}, fmt.Errorf("mock user is empty")
	}
	return "ok", time.Now().UTC().Add(ttl), nil
}

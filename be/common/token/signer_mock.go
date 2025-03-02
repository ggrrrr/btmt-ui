package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	mocker struct{}
)

var _ (Signer) = (*mocker)(nil)

func NewSignerMock() *mocker {
	log.Log().Warn(nil, "NewSignerMock")
	return &mocker{}
}

func (*mocker) Sign(_ context.Context, ttl time.Duration, claims roles.AuthInfo) (string, time.Time, error) {
	if claims.Subject == "" {
		return "", time.Time{}, fmt.Errorf("mock user is empty")
	}
	return "mockuser", time.Now().UTC().Add(ttl), nil
}

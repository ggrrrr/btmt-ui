package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	ServiceTokenGenerator interface {
		// generate token for a service name and for specific subject
		TokenForSubject(ctx context.Context, subject roles.AuthInfo) (string, error)

		// generate token for a service name
		TokenForService(ctx context.Context) (string, error)
	}

	TokenGenerator struct {
		subject string
		signer  Signer
	}
)

var _ (ServiceTokenGenerator) = (*TokenGenerator)(nil)

func NewTokenGenerator(serviceName string, signer Signer) *TokenGenerator {
	return &TokenGenerator{
		signer:  signer,
		subject: fmt.Sprintf("%s@svc", serviceName),
	}
}

// TokenForSubject implements ServiceTokenGenerator.
func (t *TokenGenerator) TokenForSubject(ctx context.Context, subject roles.AuthInfo) (string, error) {
	panic("unimplemented")
}

// TokenForService implements ServiceTokenGenerator.
func (t *TokenGenerator) TokenForService(ctx context.Context) (string, error) {
	token, _, err := t.signer.Sign(ctx, time.Hour*10000, roles.AuthInfo{
		Subject: t.subject,
	})

	if err != nil {
		return "", fmt.Errorf("TokenForSvc.signer.Sign %w", err)
	}

	return token, nil
}

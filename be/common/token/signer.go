package token

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.common.token"

type (
	signer struct {
		tracer     tracer.OTelTracer
		signMethod string
		signKey    *rsa.PrivateKey
	}

	appJwt struct {
		Roles []string `json:"roles"`
		Realm string   `json:"realm"`
		jwt.RegisteredClaims
	}
)

var _ (Signer) = (*signer)(nil)

func NewSigner(keyFile string) (*signer, error) {
	log.Log().Info("NewSigner",
		slog.String("key_file", keyFile))

	signKeyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	if err != nil {
		return nil, err
	}

	return &signer{
		tracer:     tracer.Tracer(otelScope),
		signMethod: SignMethod_RS254,
		signKey:    signKey,
	}, nil
}

func (c *signer) Sign(ctx context.Context, ttl time.Duration, authInfo roles.AuthInfo) (string, time.Time, error) {
	var err error
	_, span := c.tracer.Span(ctx, "Sign")
	defer func() {
		span.End(err)
	}()
	expiresAt := time.Now().UTC().Add(ttl)
	claims := fromAuthInfo(authInfo)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	myToken := jwt.NewWithClaims(jwt.GetSigningMethod(c.signMethod), claims)
	tokenString, err := myToken.SignedString(c.signKey)
	if err != nil {
		return "", time.Time{}, err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(tokenString))
	return encoded, expiresAt, nil
}

func fromAuthInfo(from roles.AuthInfo) *appJwt {
	out := appJwt{
		Roles: []string{},
		Realm: string(from.Realm),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: from.Subject,
			ID:      from.ID.String(),
		},
	}
	for _, v := range from.Roles {
		out.Roles = append(out.Roles, string(v))
	}

	return &out
}

package token

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	signer struct {
		ttl        time.Duration
		signMethod string
		signKey    *rsa.PrivateKey
	}

	appJwt struct {
		Roles []string `json:"roles"`
		jwt.RegisteredClaims
	}

	Signer interface {
		Sign(claims roles.AuthInfo) (string, error)
	}
)

func fromAuthInfo(from roles.AuthInfo) *appJwt {
	out := appJwt{
		Roles:            []string{},
		RegisteredClaims: jwt.RegisteredClaims{Subject: from.User},
	}
	for _, v := range from.Roles {
		out.Roles = append(out.Roles, string(v))
	}
	return &out
}

var _ (Signer) = (*signer)(nil)

func NewSigner(ttl time.Duration, keyFile string) (*signer, error) {
	signKeyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	if err != nil {
		return nil, err
	}

	logger.Log().Info().
		Str("key_file", keyFile).
		Str("ttl", ttl.String()).Send()
	return &signer{
		ttl:        ttl,
		signMethod: "RS256",
		signKey:    signKey,
	}, nil
}

func (c *signer) Sign(authInfo roles.AuthInfo) (string, error) {
	claims := fromAuthInfo(authInfo)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(c.ttl))
	mytoken := jwt.NewWithClaims(jwt.GetSigningMethod(c.signMethod), claims)
	tokenString, err := mytoken.SignedString(c.signKey)
	return tokenString, err
}

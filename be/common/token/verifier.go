package token

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

var (
	ErrJwtBadScheme         = errors.New("JWT authorization scheme is invalid")
	ErrJwtBadAlg            = errors.New("JWT Inconsistent Algorithm")
	ErrJwtInvalid           = errors.New("JWT is invalid")
	ErrJwtInvalidSubject    = errors.New("JWT subject is invalid")
	ErrJwtNotFoundMapClaims = errors.New("JWT MapClaims not found ")
)

type (
	verifier struct {
		signMethod string
		verifyKey  *rsa.PublicKey
	}

	Verifier interface {
		Verify(inputToken roles.Authorization) (roles.AuthInfo, error)
	}
)

var _ (Verifier) = (*verifier)(nil)

func NewVerifier(crtFile string) (*verifier, error) {
	crtBytes, err := os.ReadFile(crtFile)
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(crtBytes)
	if err != nil {
		return nil, err
	}
	logger.Log().Info().
		Str("crtFile", crtFile).
		Str("schema", roles.AuthSchemeBeaerer).
		Msg("Verifier")
	return &verifier{
		signMethod: "RS256",
		verifyKey:  verifyKey,
	}, nil

}

func (c *verifier) Verify(inputToken roles.Authorization) (roles.AuthInfo, error) {
	if inputToken.AuthScheme != roles.AuthSchemeBeaerer {
		return roles.AuthInfo{}, ErrJwtBadScheme

	}
	jwtToken, err := base64.StdEncoding.DecodeString(string(inputToken.AuthCredentials))
	if err != nil {
		return roles.AuthInfo{}, err
	}
	out, err := jwt.Parse(string(jwtToken), func(token *jwt.Token) (interface{}, error) {
		return c.verifyKey, nil
	})
	if err != nil {
		return roles.AuthInfo{}, err
	}
	jwtAlg := out.Header["alg"]
	if c.signMethod != jwtAlg {
		return roles.AuthInfo{}, ErrJwtBadAlg
	}
	claims, ok := out.Claims.(jwt.MapClaims)
	if !ok {
		return roles.AuthInfo{}, ErrJwtNotFoundMapClaims
	}
	if !out.Valid {
		return roles.AuthInfo{}, ErrJwtInvalid
	}

	user, ok := (claims["sub"]).(string)
	if !ok {
		return roles.AuthInfo{}, ErrJwtInvalidSubject
	}

	var listRoles []roles.RoleName
	tmp, ok := (claims["roles"]).([]interface{})
	if ok {
		listRoles = listToRoles(tmp)
	}

	return roles.AuthInfo{
		User:  user,
		Roles: listRoles,
	}, nil
}

func listToRoles(l []interface{}) []roles.RoleName {
	out := make([]roles.RoleName, 0, len(l))
	for _, r := range l {
		role, ok := r.(string)
		if ok {
			out = append(out, roles.RoleName(role))
		}
	}
	return out
}

package token

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

var (
	ErrJwtBadScheme         = errors.New("JWT authorization scheme is invalid")
	ErrJwtBadAlg            = errors.New("JWT Inconsistent Algorithm")
	ErrJwtInvalid           = errors.New("JWT is invalid")
	ErrJwtNotFoundRealm     = errors.New("JWT tenant not set")
	ErrJwtInvalidSubject    = errors.New("JWT subject is invalid")
	ErrJwtNotFoundMapClaims = errors.New("JWT MapClaims not found")
)

const (
	realmKey string = "realm"
	idKey    string = "jti"
	subKey   string = "sub"
	rolesKey string = "roles"
	algKey   string = "alg"
)

type (
	verifier struct {
		signMethod string
		verifyKey  *rsa.PublicKey
	}

	Verifier interface {
		Verify(inputToken app.AuthData) (roles.AuthInfo, error)
	}
)

var _ (Verifier) = (*verifier)(nil)

func NewVerifier(crtFile string) (*verifier, error) {
	logger.Info().
		Str("crtFile", crtFile).
		Str("schema", roles.AuthSchemeBearer).
		Msg("NewVerifier")

	crtBytes, err := os.ReadFile(crtFile)
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(crtBytes)
	if err != nil {
		return nil, err
	}

	return &verifier{
		signMethod: "RS256",
		verifyKey:  verifyKey,
	}, nil

}

func (c *verifier) Verify(inputToken app.AuthData) (roles.AuthInfo, error) {
	if inputToken.AuthScheme != roles.AuthSchemeBearer {
		return roles.AuthInfo{}, ErrJwtBadScheme

	}

	jwtTokenBytes, err := base64.StdEncoding.DecodeString(string(inputToken.AuthToken))
	if err != nil {
		return roles.AuthInfo{}, err
	}

	jwtToken, err := jwt.Parse(string(jwtTokenBytes), func(token *jwt.Token) (interface{}, error) {
		return c.verifyKey, nil
	})
	if err != nil {
		return roles.AuthInfo{}, err
	}

	jwtAlg := jwtToken.Header[algKey]
	if c.signMethod != jwtAlg {
		return roles.AuthInfo{}, ErrJwtBadAlg
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return roles.AuthInfo{}, ErrJwtNotFoundMapClaims
	}

	if !jwtToken.Valid {
		return roles.AuthInfo{}, ErrJwtInvalid
	}

	subject, ok := (claims[subKey]).(string)
	if !ok {
		return roles.AuthInfo{}, ErrJwtInvalidSubject
	}

	id, _ := (claims[idKey]).(string)

	var realm string
	var listRoles []string
	tmp, ok := (claims[rolesKey]).([]interface{})
	if ok {
		listRoles = listToRoles(tmp)
	}

	realm, ok = (claims[realmKey]).(string)
	if !ok {
		return roles.AuthInfo{}, ErrJwtNotFoundRealm
	}

	if realm == "" {
		return roles.AuthInfo{}, ErrJwtNotFoundRealm
	}

	authInfo := roles.AuthInfo{
		Subject: subject,
		Realm:   realm,
		Roles:   listRoles,
	}

	tokenID, err := uuid.Parse(id)
	if err == nil {
		authInfo.ID = tokenID
	}
	return authInfo, nil
}

func listToRoles(l []interface{}) []string {
	out := make([]string, 0, len(l))
	for _, r := range l {
		role, ok := r.(string)
		if ok {
			out = append(out, role)
		}
	}
	return out
}

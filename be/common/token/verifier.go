package token

import (
	"crypto/rsa"
	"encoding/base64"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	verifier struct {
		// tracer     tracer.OTelTracer
		signMethod string
		verifyKey  *rsa.PublicKey
	}
)

var _ (Verifier) = (*verifier)(nil)

func NewVerifier(crtFile string) (*verifier, error) {
	log.Log().Info("NewVerifier",
		slog.String("crtFile", crtFile),
		slog.String("schema", roles.AuthSchemeBearer))

	crtBytes, err := os.ReadFile(crtFile)
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(crtBytes)
	if err != nil {
		return nil, err
	}

	return &verifier{
		// tracer:     tracer.Tracer(otelScope),
		signMethod: SignMethod_RS254,
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

	jwtAlg := jwtToken.Header[jwtAlgorithmKey]
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

	subject, ok := (claims[jwtSubjectKey]).(string)
	if !ok {
		return roles.AuthInfo{}, ErrJwtInvalidSubject
	}

	id, _ := (claims[jwtIDKey]).(string)

	var realm string
	var listRoles []string
	tmp, ok := (claims[jwtRolesKey]).([]interface{})
	if ok {
		listRoles = listToRoles(tmp)
	}

	realm, ok = (claims[jwtRealmKey]).(string)
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

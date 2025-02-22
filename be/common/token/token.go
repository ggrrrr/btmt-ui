package token

import (
	"context"
	"errors"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

const SignMethod_RS254 string = "RS256"

var (
	ErrJwtBadScheme         = errors.New("JWT authorization scheme is invalid")
	ErrJwtBadAlg            = errors.New("JWT Inconsistent Algorithm")
	ErrJwtInvalid           = errors.New("JWT is invalid")
	ErrJwtNotFoundRealm     = errors.New("JWT tenant not set")
	ErrJwtInvalidSubject    = errors.New("JWT subject is invalid")
	ErrJwtNotFoundMapClaims = errors.New("JWT MapClaims not found")
)

type (
	Config struct {
		CrtFile string `env:"CRT_FILE"`
		KeyFile string `env:"KEY_FILE"`
	}

	Verifier interface {
		Verify(inputToken app.AuthData) (roles.AuthInfo, error)
	}

	Signer interface {
		Sign(ctx context.Context, ttl time.Duration, claims roles.AuthInfo) (string, time.Time, error)
	}
)

const (
	jwtRealmKey     string = "realm"
	jwtIDKey        string = "jti"
	jwtSubjectKey   string = "sub"
	jwtRolesKey     string = "roles"
	jwtAlgorithmKey string = "alg"
)

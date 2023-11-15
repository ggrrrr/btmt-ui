package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func TestSignTTL(t *testing.T) {
	testSigner, err := NewSigner(1*time.Second, "/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/be/jwt.key")
	assert.NoError(t, err)

	testVer, err := NewVerifier("/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/be/jwt.crt")
	assert.NoError(t, err)

	apiClaims := roles.AuthInfo{
		User:  "user1",
		Roles: []roles.RoleName{},
	}

	jwt, err := testSigner.Sign(apiClaims)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	_, err = testVer.Verify(
		roles.Authorization{
			AuthScheme:      "",
			AuthCredentials: roles.AuthCredentials(jwt),
		})
	logger.Log().Info().Any("err", err).Msg("v")
	assert.Error(t, err)
}

func TestSignVerify(t *testing.T) {
	testSigner, err := NewSigner(1*time.Second, "/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/be/jwt.key")
	assert.NoError(t, err)

	testVer, err := NewVerifier("/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/be/jwt.crt")
	assert.NoError(t, err)

	apiClaims := roles.AuthInfo{
		User:  "user1",
		Roles: []roles.RoleName{"admin"},
	}

	expClaims := roles.AuthInfo{
		User:  "user1",
		Roles: []roles.RoleName{"admin"},
	}

	_, err = testVer.Verify(
		roles.Authorization{
			AuthScheme:      "asdasd",
			AuthCredentials: roles.AuthCredentials("asdasdasd"),
		},
	)
	assert.ErrorIs(t, err, ErrJwtBadScheme)

	jwt, err := testSigner.Sign(apiClaims)
	assert.NoError(t, err)
	c, err := testVer.Verify(
		roles.Authorization{
			AuthScheme:      roles.AuthSchemeBeaerer,
			AuthCredentials: roles.AuthCredentials(jwt),
		},
	)
	assert.NoError(t, err)
	logger.Log().Info().Any("c", c).Msg("v")
	assert.Equal(t, expClaims, c)
}

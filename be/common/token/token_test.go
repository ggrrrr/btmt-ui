package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/help"
)

func TestSignTTL(t *testing.T) {
	pwd := help.RepoDir()

	testSigner, err := NewSigner(1*time.Second, fmt.Sprintf("%s/jwt.key", pwd))
	require.NoError(t, err)

	testVer, err := NewVerifier(fmt.Sprintf("%s/jwt.crt", pwd))
	require.NoError(t, err)

	apiClaims := roles.AuthInfo{
		User:  "user1",
		Roles: []string{},
	}

	jwt, err := testSigner.Sign(apiClaims)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	_, err = testVer.Verify(
		roles.Authorization{
			AuthScheme:      "",
			AuthCredentials: roles.AuthCredentials(jwt),
		})
	logger.Info().Any("err", err).Msg("v")
	assert.Error(t, err)
}

func TestSignVerify(t *testing.T) {
	testSigner, err := NewSigner(1*time.Second, "/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/jwt.key")
	require.NoError(t, err)

	testVer, err := NewVerifier("/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/jwt.crt")
	require.NoError(t, err)

	apiClaims := roles.AuthInfo{
		User:   "user1",
		Tenant: "localhost",
		Roles:  []string{"admin"},
	}

	expClaims := roles.AuthInfo{
		User:   "user1",
		Tenant: "localhost",
		Roles:  []string{"admin"},
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
			AuthScheme:      roles.AuthSchemeBearer,
			AuthCredentials: roles.AuthCredentials(jwt),
		},
	)
	assert.NoError(t, err)
	logger.Info().Any("c", c).Msg("v")
	assert.Equal(t, expClaims, c)
}

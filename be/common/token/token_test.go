package token

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/help"
)

func TestSignTTL(t *testing.T) {
	ctx := context.TODO()

	ttl := 1 * time.Second
	pwd := help.RepoDir()

	testSigner, err := NewSigner(fmt.Sprintf("%s/jwt.key", pwd))
	require.NoError(t, err)

	testVer, err := NewVerifier(fmt.Sprintf("%s/jwt.crt", pwd))
	require.NoError(t, err)

	apiClaims := roles.AuthInfo{
		Subject: "user1",
		Roles:   []string{},
	}

	jwt, expiresAt, err := testSigner.Sign(ctx, ttl, apiClaims)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	_, err = testVer.Verify(
		app.AuthData{
			AuthScheme: "",
			AuthToken:  jwt,
		})
	logger.Info().Any("err", err).Msg("v")
	assert.Error(t, err)
	assert.True(t, !expiresAt.IsZero())
}

func TestSignVerify(t *testing.T) {
	ctx := context.TODO()
	pwd := help.RepoDir()

	ttl := 1 * time.Second
	testSigner, err := NewSigner(fmt.Sprintf("%s/jwt.key", pwd))
	require.NoError(t, err)

	testVer, err := NewVerifier(fmt.Sprintf("%s/jwt.crt", pwd))
	require.NoError(t, err)

	tokenID := uuid.New()

	apiClaims := roles.AuthInfo{
		Subject: "user1",
		Realm:   "localhost",
		Roles:   []string{"admin"},
		ID:      tokenID,
	}

	expClaims := roles.AuthInfo{
		Subject: "user1",
		Realm:   "localhost",
		Roles:   []string{"admin"},
		ID:      tokenID,
	}

	_, err = testVer.Verify(
		app.AuthData{
			AuthScheme: "asdasd",
			AuthToken:  "asdasdasd",
		},
	)
	assert.ErrorIs(t, err, ErrJwtBadScheme)

	jwt, expiresAt, err := testSigner.Sign(ctx, ttl, apiClaims)
	assert.NoError(t, err)
	authInfo, err := testVer.Verify(
		app.AuthData{
			AuthScheme: roles.AuthSchemeBearer,
			AuthToken:  jwt,
		},
	)
	assert.NoError(t, err)
	logger.Info().Any("c", authInfo).Msg("v")
	assert.Equal(t, expClaims, authInfo)
	assert.True(t, !expiresAt.IsZero())
}

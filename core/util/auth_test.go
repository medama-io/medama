package util_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupAuthTest(t *testing.T) (*assert.Assertions, *require.Assertions, context.Context, *util.AuthService) {
	t.Helper()
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()

	auth, err := util.NewAuthService(ctx, false)
	require.NoError(err)
	assert.NotNil(auth)

	return assert, require, ctx, auth
}

func TestAuthCreateAndRead(t *testing.T) {
	assert, require, ctx, auth := SetupAuthTest(t)

	// We don't want the cookie to expire
	cookie, err := auth.CreateSession(ctx, "test_user_id")
	require.NoError(err)
	assert.NotNil(cookie)

	assert.Equal("_me_sess", cookie.Name)
	assert.Equal("/", cookie.Path)
	assert.True(cookie.HttpOnly)
	assert.True(cookie.Secure)
	assert.Equal(http.SameSiteLaxMode, cookie.SameSite)

	// Decrypt cookie
	userId, err := auth.ReadSession(ctx, cookie.Value)
	require.NoError(err)
	assert.Equal("test_user_id", userId)
}

func TestAuthWithInvalidSession(t *testing.T) {
	assert, require, ctx, auth := SetupAuthTest(t)

	// We don't want the cookie to expire
	cookie, err := auth.CreateSession(ctx, "test_user")
	require.NoError(err)
	assert.NotNil(cookie)

	// Decrypt cookie
	userId, err := auth.ReadSession(ctx, "invalid_session")
	require.ErrorIs(err, model.ErrInvalidSession)
	assert.Equal("", userId)
}

func TestAuthWithExpiredSession(t *testing.T) {
	assert, require, ctx, auth := SetupAuthTest(t)

	cookie, err := auth.CreateSession(ctx, "test_user_id")
	require.NoError(err)
	assert.NotNil(cookie)

	// Delete from cache to simulate expired session
	base64Decode, err := base64.URLEncoding.DecodeString(cookie.Value)
	require.NoError(err)
	sessionId, err := auth.DecryptSession(ctx, string(base64Decode))
	require.NoError(err)
	assert.NotEmpty(sessionId)
	auth.Cache.Delete(sessionId)

	// Try to read from session with expired cookie
	userId, err := auth.ReadSession(ctx, cookie.Value)
	require.ErrorIs(err, model.ErrSessionNotFound)
	assert.Equal("", userId)
}

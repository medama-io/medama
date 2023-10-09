package util_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"github.com/stretchr/testify/assert"
)

func SetupAuthTest(t *testing.T) (*assert.Assertions, context.Context, *util.AuthService) {
	t.Helper()
	assert := assert.New(t)
	ctx := context.Background()

	auth, err := util.NewAuthService(ctx)
	assert.NoError(err)
	assert.NotNil(auth)

	return assert, ctx, auth
}

func TestAuthEncryptAndDecrypt(t *testing.T) {
	assert, ctx, auth := SetupAuthTest(t)

	// We don't want the cookie to expire
	cookie, err := auth.CreateSession(ctx, "test_user", 24*time.Hour)
	assert.NoError(err)
	assert.NotNil(cookie)

	assert.Equal("_me_sess", cookie.Name)
	assert.Equal("/", cookie.Path)
	assert.True(cookie.HttpOnly)
	assert.True(cookie.Secure)
	assert.Equal(http.SameSiteStrictMode, cookie.SameSite)

	// Decrypt cookie
	userId, err := auth.ReadSession(ctx, cookie.Value)
	assert.NoError(err)
	assert.Equal("test_user", userId)
}

func TestAuthWithInvalidSession(t *testing.T) {
	assert, ctx, auth := SetupAuthTest(t)

	// We don't want the cookie to expire
	cookie, err := auth.CreateSession(ctx, "test_user", 24*time.Hour)
	assert.NoError(err)
	assert.NotNil(cookie)

	// Decrypt cookie
	userId, err := auth.ReadSession(ctx, "invalid_session")
	assert.ErrorIs(err, model.ErrInvalidSession)
	assert.Equal("", userId)
}

func TestAuthWithExpiredSession(t *testing.T) {
	assert, ctx, auth := SetupAuthTest(t)

	cookie, err := auth.CreateSession(ctx, "test_user", 24*time.Hour)
	assert.NoError(err)
	assert.NotNil(cookie)

	// Delete from cache to simulate expired session
	sessionId, err := auth.DecryptSession(ctx, cookie.Value)
	assert.NoError(err)
	assert.NotEmpty(sessionId)
	auth.Cache.Delete(sessionId)

	// Try to read from session with expired cookie
	userId, err := auth.ReadSession(ctx, cookie.Value)
	assert.ErrorIs(err, model.ErrSessionNotFound)
	assert.Equal("", userId)
}

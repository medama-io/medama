package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/metest"
	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

var ValidApiKey string = "valid_api_key"

func TestGetsUserApiKeyWhenExists(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)
	ctx = context.WithValue(ctx, model.ContextKeyUserID, "api_user_1")

	user := model.NewUser(
		"api_user_1",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	user.ApiKey = &ValidApiKey

	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.GetUserAPIKey(ctx, api.GetUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.UserApiKeyHeaders{}, res)

	responseHeaders := res.(*api.UserApiKeyHeaders)
	assert.Equal(responseHeaders.Response.APIKey.Value, ValidApiKey)
}

func TestGetUserApiKeyWhenKeyDoesNotExists(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)
	ctx = context.WithValue(ctx, model.ContextKeyUserID, "api_user_2")

	user := model.NewUser(
		"api_user_2",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.GetUserAPIKey(ctx, api.GetUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.NotFoundErrorHeaders{}, res)
}

func TestGetUserApiKeyWhenUnauthorized(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)

	user := model.NewUser(
		"api_user_3",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.GetUserAPIKey(ctx, api.GetUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.UnauthorisedErrorHeaders{}, res)
}

func TestGetUserApiKeyWhenUserNotExists(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)
	ctx = context.WithValue(ctx, model.ContextKeyUserID, "api_user_1")

	user := model.NewUser(
		"api_user_4",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.GetUserAPIKey(ctx, api.GetUserAPIKeyParams{})

	assert.Nil(res)
	assert.ErrorIs(err, model.ErrUserNotFound)
}

func TestRegeneratesUserApiKey(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)
	ctx = context.WithValue(ctx, model.ContextKeyUserID, "api_user_5")

	user := model.NewUser(
		"api_user_5",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	user.ApiKey = &ValidApiKey

	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.RegenerateUserAPIKey(ctx, api.RegenerateUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.UserApiKeyHeaders{}, res)

	responseHeaders := res.(*api.UserApiKeyHeaders)

	dbUser, err := database.GetUser(ctx, "api_user_5")

	require.NoError(t, err)
	assert.NotNil(dbUser)
	assert.Equal(responseHeaders.Response.APIKey.Value, *dbUser.ApiKey)
}

func TestRegenerateApiKeyReturnsUnauthorized(t *testing.T) {
	assert, ctx, handler, _ := metest.NewTestHandler(t)

	res, err := handler.RegenerateUserAPIKey(ctx, api.RegenerateUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.UnauthorisedErrorHeaders{}, res)
}

func TestRevokesUserApiKey(t *testing.T) {
	assert, ctx, handler, database := metest.NewTestHandler(t)
	ctx = context.WithValue(ctx, model.ContextKeyUserID, "api_user_6")

	user := model.NewUser(
		"api_user_6",
		"api_user",
		"",
		model.NewDefaultUserSettings(),
		time.Now().Unix(),
		time.Now().Unix(),
	)
	user.ApiKey = &ValidApiKey

	err := database.CreateUser(ctx, user)
	require.NoError(t, err)

	res, err := handler.DeleteUserAPIKey(ctx, api.DeleteUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.DeleteUserAPIKeyNoContent{}, res)

	dbUser, err := database.GetUser(ctx, "api_user_6")

	require.NoError(t, err)
	assert.NotNil(dbUser)
	assert.Nil(dbUser.ApiKey)
}

func TestRevokeApiKeyReturnsUnauthorized(t *testing.T) {
	assert, ctx, handler, _ := metest.NewTestHandler(t)

	res, err := handler.DeleteUserAPIKey(ctx, api.DeleteUserAPIKeyParams{})

	require.NoError(t, err)
	assert.IsType(&api.UnauthorisedErrorHeaders{}, res)
}

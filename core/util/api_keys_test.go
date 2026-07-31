package util_test

import (
	"context"
	"testing"
	"time"

	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/metest"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"github.com/stretchr/testify/assert"
)

var ValidApiKey string = "valid_api_key"

func TestExchangeApiKeyReturnsUser(t *testing.T) {
	service, database, assert, ctx := setupApiKeysService(t)

	user := model.NewUser("api_user_1", "api_user", "", model.NewDefaultUserSettings(), time.Now().Unix(), time.Now().Unix())
	user.ApiKey = &ValidApiKey

	database.CreateUser(ctx, user)

	user, err := service.ExchangeApiKey(ctx, ValidApiKey)

	assert.NoError(err)
	assert.IsType(&model.User{}, user)
	assert.Equal("api_user_1", user.ID)
}

func TestExchangeApiKeyReturnsInvalidKey(t *testing.T) {
	service, _, assert, context := setupApiKeysService(t)

	user, err := service.ExchangeApiKey(context, "invalid_api_key")

	assert.Nil(user)
	assert.ErrorIs(err, model.ErrInvalidApiKey)
}

func TestRegenerateUserApiKey(t *testing.T) {
	service, database, assert, ctx := setupApiKeysService(t)

	user := model.NewUser("api_user_2", "api_user", "", model.NewDefaultUserSettings(), time.Now().Unix(), time.Now().Unix())
	user.ApiKey = &ValidApiKey

	database.CreateUser(ctx, user)

	apiKey, err := service.RegenerateUserApiKey(ctx, "api_user_2")

	assert.NoError(err)

	dbUser, err := database.GetUser(ctx, "api_user_2")

	assert.NoError(err)
	assert.Equal(apiKey, *dbUser.ApiKey)
}

func TestRevokesUserApiKey(t *testing.T) {
	service, database, assert, ctx := setupApiKeysService(t)

	user := model.NewUser("api_user_3", "api_user", "", model.NewDefaultUserSettings(), time.Now().Unix(), time.Now().Unix())
	user.ApiKey = &ValidApiKey

	database.CreateUser(ctx, user)

	err := service.RevokeUserApiKey(ctx, "api_user_3")

	assert.NoError(err)

	dbUser, err := database.GetUser(ctx, "api_user_3")

	assert.NoError(err)
	assert.NotNil(dbUser)
	assert.Nil(dbUser.ApiKey)
}

func setupApiKeysService(t *testing.T) (*util.ApiKeysService, *sqlite.Client, *assert.Assertions, context.Context) {
	assert := assert.New(t)
	ctx := t.Context()

	appDatabase, _ := metest.NewInMemoryDatabase(t)
	service, err := util.NewApiKeysService(appDatabase)

	assert.NoError(err)
	assert.NotNil(service)

	return service, appDatabase, assert, ctx
}

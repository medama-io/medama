package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	userCreate := model.NewUser(
		"test",
		"username",
		"password",
		model.NewDefaultSettings(),
		1,
		2,
	)

	err := client.CreateUser(ctx, userCreate)
	require.NoError(t, err)

	user, err := client.GetUser(ctx, "test")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test", user.ID)
	assert.Equal("username", user.Username)
	assert.Equal("en", user.Settings.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)
	assert.Equal("en", user.Settings.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetUserNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUser(ctx, "doesnotexist")
	require.ErrorIs(t, err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestGetUserByUsername(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUserByUsername(ctx, "username1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)
	assert.Equal("en", user.Settings.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetDefaultAdminUser(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "admin")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("admin", user.Username)
	assert.Equal("en", user.Settings.Language)
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "doesnotexist")
	require.ErrorIs(t, err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestUpdateUserUsername(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)

	dateUpdated := user.DateUpdated

	err = client.UpdateUserUsername(ctx, "test1", "usernamenew")
	require.NoError(t, err)

	user, err = client.GetUser(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("usernamenew", user.Username)
	assert.Greater(user.DateUpdated, dateUpdated)
}

func TestUpdateUserUsernameExisting(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)

	err = client.UpdateUserUsername(ctx, "test1", "username2")
	assert.ErrorIs(err, model.ErrUserExists)
}

func TestUpdateUserPassword(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	err := client.UpdateUserPassword(ctx, "test1", "password2")
	require.NoError(t, err)

	user, err := client.GetUser(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("password2", user.Password)
}

func TestDeleteUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	err := client.DeleteUser(ctx, "test1")
	require.NoError(t, err)

	user, err := client.GetUser(ctx, "test1")
	require.ErrorIs(t, err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestDeleteUserNotFound(t *testing.T) {
	_, ctx, client := SetupDatabase(t)

	err := client.DeleteUser(ctx, "doesnotexist")
	require.ErrorIs(t, err, model.ErrUserNotFound)
}

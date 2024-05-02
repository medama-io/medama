package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestCreateUser(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	userCreate := model.NewUser(
		"test",
		"username",
		"password",
		"en",
		1,
		2,
	)

	err := client.CreateUser(ctx, userCreate)
	assert.NoError(err)

	user, err := client.GetUser(ctx, "test")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test", user.ID)
	assert.Equal("username", user.Username)
	assert.Equal("en", user.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)
	assert.Equal("en", user.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetUserNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUser(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestGetUserByUsername(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUserByUsername(ctx, "username1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)
	assert.Equal("en", user.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetDefaultAdminUser(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "admin")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("admin", user.Username)
	assert.Equal("en", user.Language)
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestUpdateUserUsername(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)

	err = client.UpdateUserUsername(ctx, "test1", "usernamenew", 3)
	assert.NoError(err)

	user, err = client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("usernamenew", user.Username)
	assert.Equal(int64(3), user.DateUpdated)
}

func TestUpdateUserUsernameExisting(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("username1", user.Username)

	err = client.UpdateUserUsername(ctx, "test1", "username2", 3)
	assert.ErrorIs(err, model.ErrUserExists)
}

func TestUpdateUserPassword(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	err := client.UpdateUserPassword(ctx, "test1", "password2", 3)
	assert.NoError(err)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("password2", user.Password)
	assert.Equal(int64(3), user.DateUpdated)
}

func TestDeleteUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	err := client.DeleteUser(ctx, "test1")
	assert.NoError(err)

	user, err := client.GetUser(ctx, "test1")
	assert.ErrorIs(err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestDeleteUserNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.DeleteUser(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrUserNotFound)
}

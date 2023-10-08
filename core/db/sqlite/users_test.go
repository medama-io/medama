package sqlite_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/medama-io/medama/model"
)

func TestCreateUser(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	userCreate := &model.User{
		ID:          "test",
		Email:       "email@example.com",
		Password:    "password",
		Language:    "en",
		DateCreated: 1,
		DateUpdated: 2,
	}

	err := client.CreateUser(ctx, userCreate)
	assert.NoError(err)

	user, err := client.GetUser(ctx, "test")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test", user.ID)
	assert.Equal("email@example.com", user.Email)
	assert.Equal("en", user.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestCreateUserDuplicate(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	userCreate := &model.User{
		Password:    "password",
		Language:    "en",
		DateCreated: 1,
		DateUpdated: 2,
	}

	// Test unique id
	userCreate.ID = "test"
	userCreate.Email = "id@example.com"
	err := client.CreateUser(ctx, userCreate)
	assert.NoError(err)

	// Should give a duplicate error for id
	userCreate.Email = "id2@example.com"
	err = client.CreateUser(ctx, userCreate)
	assert.ErrorIs(err, model.ErrUserExists)

	// Test unique email
	userCreate.ID = "test2"
	userCreate.Email = "email@example.com"
	err = client.CreateUser(ctx, userCreate)
	assert.NoError(err)

	// Should give a duplicate error for email
	userCreate.ID = "test3"
	err = client.CreateUser(ctx, userCreate)
	assert.ErrorIs(err, model.ErrUserExists)
}

func TestGetUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("test1@example.com", user.Email)
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

func TestGetUserByEmail(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUserByEmail(ctx, "test1@example.com")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("test1@example.com", user.Email)
	assert.Equal("en", user.Language)
	assert.Equal(int64(1), user.DateCreated)
	assert.Equal(int64(2), user.DateUpdated)
}

func TestGetUserByEmailNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByEmail(ctx, "doesnotexist@example.com")
	assert.ErrorIs(err, model.ErrUserNotFound)
	assert.Nil(user)
}

func TestUpdateUserEmail(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("test1@example.com", user.Email)

	err = client.UpdateUserEmail(ctx, "test1", "testUpdate1@example.com", 3)
	assert.NoError(err)

	user, err = client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
	assert.Equal("testUpdate1@example.com", user.Email)
	assert.Equal(int64(3), user.DateUpdated)
}

func TestUpdateUserPassword(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	err := client.UpdateUserPassword(ctx, "test1", "password2", 3)
	assert.NoError(err)

	// Custom query since password is not returned by GetUser
	var password string
	query := `--sql
	SELECT password FROM users WHERE id = ?`
	err = client.DB.QueryRowxContext(ctx, query, "test1").Scan(&password)
	assert.NoError(err)
	assert.Equal("password2", password)

	user, err := client.GetUser(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("test1", user.ID)
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

package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

func TestCreateWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := model.NewWebsite(
		"test1",
		"example.com",
		1,
		2,
	)

	err := client.CreateWebsite(ctx, websiteCreate)
	require.NoError(t, err)

	website, err := client.GetWebsite(ctx, "example.com")
	require.NoError(t, err)
	assert.NotNil(website)
	assert.Equal("test1", website.UserID)
	assert.Equal("example.com", website.Hostname)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestCreateWebsiteDuplicateHostname(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := model.NewWebsite(
		"test1",
		"example.com",
		1,
		2,
	)

	// Test unique id
	err := client.CreateWebsite(ctx, websiteCreate)
	require.NoError(t, err)

	// Should give a duplicate error for id
	err = client.CreateWebsite(ctx, websiteCreate)
	require.ErrorIs(t, err, model.ErrWebsiteExists)

	// Test unique email
	websiteCreate.Hostname = "example2.com"
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)
}

func TestCreateWebsiteMissingUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := model.NewWebsite(
		"doesnotexist",
		"exampledoesnotexist.com",
		1,
		2,
	)

	err := client.CreateWebsite(ctx, websiteCreate)
	assert.ErrorIs(err, model.ErrUserNotFound)
}

func TestListWebsites(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	websites, err := client.ListWebsites(ctx, "test1")
	require.NoError(t, err)
	assert.NotNil(websites)
	assert.Len(websites, 3)
}

func TestListWebsitesNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	websites, err := client.ListWebsites(ctx, "doesnotexist")
	require.NoError(t, err)
	assert.NotNil(websites)
	assert.Empty(websites)
}

func TestGetWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	website, err := client.GetWebsite(ctx, "website1-test1.com")
	require.NoError(t, err)
	assert.NotNil(website)
	assert.Equal("test1", website.UserID)
	assert.Equal("website1-test1.com", website.Hostname)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestGetWebsiteNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	website, err := client.GetWebsite(ctx, "doesnotexist.com")
	require.ErrorIs(t, err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestListAllHostnames(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	hostnames, err := client.ListAllHostnames(ctx)
	require.NoError(t, err)
	assert.NotNil(hostnames)
	// 3 websites each for 3 users
	assert.Len(hostnames, 9)
}

func TestDeleteWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	website, err := client.GetWebsite(ctx, "website1-test1.com")
	require.NoError(t, err)
	assert.NotNil(website)

	err = client.DeleteWebsite(ctx, "website1-test1.com")
	require.NoError(t, err)

	website, err = client.GetWebsite(ctx, "website1-test1.com")
	require.ErrorIs(t, err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestDeleteWebsiteNotFound(t *testing.T) {
	_, ctx, client := SetupDatabase(t)

	err := client.DeleteWebsite(ctx, "doesnotexist")
	require.ErrorIs(t, err, model.ErrWebsiteNotFound)
}

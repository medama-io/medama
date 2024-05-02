package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestCreateWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := model.NewWebsite(
		"test1",
		"example.com",
		"Example",
		1,
		2,
	)

	err := client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	website, err := client.GetWebsite(ctx, "example.com")
	assert.NoError(err)
	assert.NotNil(website)
	assert.Equal("test1", website.UserID)
	assert.Equal("example.com", website.Hostname)
	assert.Equal("Example", website.Name)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestCreateWebsiteDuplicateHostname(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := model.NewWebsite(
		"test1",
		"example.com",
		"Example",
		1,
		2,
	)

	// Test unique id
	err := client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	// Should give a duplicate error for id
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.Error(err, model.ErrWebsiteExists)

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
		"DoesNotExist",
		1,
		2,
	)

	err := client.CreateWebsite(ctx, websiteCreate)
	assert.ErrorIs(err, model.ErrUserNotFound)
}

func TestListWebsites(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	websites, err := client.ListWebsites(ctx, "test1")
	assert.NoError(err)
	assert.NotNil(websites)
	assert.Len(websites, 3)
}

func TestListWebsitesNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	websites, err := client.ListWebsites(ctx, "doesnotexist")
	assert.Nil(websites)
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
}

func TestGetWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	website, err := client.GetWebsite(ctx, "website1-test1.com")
	assert.NoError(err)
	assert.NotNil(website)
	assert.Equal("test1", website.UserID)
	assert.Equal("website1-test1.com", website.Hostname)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestGetWebsiteNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	website, err := client.GetWebsite(ctx, "doesnotexist.com")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestWebsiteExists(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	exists, err := client.WebsiteExists(ctx, "website1-test1.com")
	assert.NoError(err)
	assert.True(exists)
}

func TestWebsiteExistsNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	exists, err := client.WebsiteExists(ctx, "doesnotexist.com")
	assert.NoError(err)
	assert.False(exists)
}

func TestDeleteWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	website, err := client.GetWebsite(ctx, "website1-test1.com")
	assert.NoError(err)
	assert.NotNil(website)

	err = client.DeleteWebsite(ctx, "website1-test1.com")
	assert.NoError(err)

	website, err = client.GetWebsite(ctx, "website1-test1.com")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestDeleteWebsiteNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.DeleteWebsite(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
}

package sqlite_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/medama-io/medama/model"
)

func TestCreateWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := &model.Website{
		ID:          "test",
		UserID:      "test1",
		Hostname:    "example.com",
		DateCreated: 1,
		DateUpdated: 2,
	}

	err := client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	website, err := client.GetWebsite(ctx, "test")
	assert.NoError(err)
	assert.NotNil(website)
	assert.Equal("test", website.ID)
	assert.Equal("test1", website.UserID)
	assert.Equal("example.com", website.Hostname)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestCreateWebsiteDuplicateID(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := &model.Website{
		ID:          "test",
		UserID:      "test1",
		Hostname:    "example.com",
		DateCreated: 1,
		DateUpdated: 2,
	}

	// Test unique id
	err := client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	websiteCreate.Hostname = "example2.com"
	// Should give a duplicate error for id
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.Error(err, model.ErrWebsiteExists)
}

func TestCreateWebsiteDuplicateHostname(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := &model.Website{
		UserID:      "test1",
		Hostname:    "example.com",
		DateCreated: 1,
		DateUpdated: 2,
	}

	// Test unique id
	websiteCreate.ID = "test"
	err := client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	// Should give a duplicate error for id
	websiteCreate.ID = "test2"
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.Error(err, model.ErrWebsiteExists)

	// Test unique email
	websiteCreate.ID = "test3"
	websiteCreate.Hostname = "example2.com"
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)
}

func TestCreateWebsiteMissingUser(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	websiteCreate := &model.Website{
		ID:          "test",
		UserID:      "doesnotexist",
		Hostname:    "exampledoesnotexist.com",
		DateCreated: 1,
		DateUpdated: 2,
	}

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

	website, err := client.GetWebsite(ctx, "website1-test1")
	assert.NoError(err)
	assert.NotNil(website)
	assert.Equal("website1-test1", website.ID)
	assert.Equal("test1", website.UserID)
	assert.Equal("website1-test1.com", website.Hostname)
	assert.Equal(int64(1), website.DateCreated)
	assert.Equal(int64(2), website.DateUpdated)
}

func TestGetWebsiteNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	website, err := client.GetWebsite(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestDeleteWebsite(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithWebsites(t)

	website, err := client.GetWebsite(ctx, "website1-test1")
	assert.NoError(err)
	assert.NotNil(website)

	err = client.DeleteWebsite(ctx, "website1-test1")
	assert.NoError(err)

	website, err = client.GetWebsite(ctx, "website1-test1")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
	assert.Nil(website)
}

func TestDeleteWebsiteNotFound(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.DeleteWebsite(ctx, "doesnotexist")
	assert.ErrorIs(err, model.ErrWebsiteNotFound)
}

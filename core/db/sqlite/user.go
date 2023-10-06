package sqlite

import "github.com/medama-io/medama/model"

func (c *Client) CreateUser(user *model.User) error {
	return nil
}

func (c *Client) GetUser(id string) (*model.User, error) {
	return nil, nil
}

func (c *Client) GetUserByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (c *Client) UpdateUser(user *model.User) error {
	return nil
}

func (c *Client) DeleteUser(id string) error {
	return nil
}

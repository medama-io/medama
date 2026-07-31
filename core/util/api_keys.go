package util

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

const (
	ApiKeyLength = 64
)

type ApiKeysService struct {
	client db.AppClient
}

func NewApiKeysService(
	client db.AppClient,
) (*ApiKeysService, error) {
	return &ApiKeysService{
		client: client,
	}, nil
}

func (s *ApiKeysService) ExchangeApiKey(ctx context.Context, apiKey string) (*model.User, error) {
	user, err := s.client.GetUserByApiKey(ctx, apiKey)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidApiKey
		}

		return nil, errors.Wrap(err, "api keys service")
	}

	return user, nil
}

func (s *ApiKeysService) RegenerateUserApiKey(ctx context.Context, userID string) (string, error) {
	token, err := GenerateRandomString(ApiKeyLength)
	if err != nil {
		return "", errors.Wrap(err, "api keys service")
	}

	err = s.client.UpdateUserApiKey(ctx, userID, &token)
	if err != nil {
		return "", errors.Wrap(err, "api keys service")
	}

	return token, nil
}

func (s *ApiKeysService) RevokeUserApiKey(ctx context.Context, userID string) error {
	err := s.client.UpdateUserApiKey(ctx, userID, nil)
	if err != nil {
		return errors.Wrap(err, "api keys service")
	}

	return nil
}

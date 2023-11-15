package access_token

import (
	"github.com/loxt/bookstore-oauth-api/internal/utils/errors"
	"strings"
)

type Repository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type Service interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)

	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	return s.repository.GetById(accessTokenId)
}

package db

import (
	"github.com/loxt/bookstore-oauth-api/internal/domain/access_token"
	"github.com/loxt/bookstore-oauth-api/internal/utils/errors"
)

type Repository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
}

type repository struct {
}

func New() Repository {
	return &repository{}
}

func (r *repository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("database connection not implemented yet")
}

package db

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/loxt/bookstore-oauth-api/internal/client/cassandra"
	"github.com/loxt/bookstore-oauth-api/internal/domain/access_token"
	errorUtils "github.com/loxt/bookstore-oauth-api/internal/utils/errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires from access_tokens WHERE access_token = ?;"
	queryCreateAccessToken = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type Repository interface {
	GetById(string) (*access_token.AccessToken, *errorUtils.RestErr)
	Create(token access_token.AccessToken) *errorUtils.RestErr
	UpdateExpirationTime(token access_token.AccessToken) *errorUtils.RestErr
}

type repository struct {
}

func New() Repository {
	return &repository{}
}

func (r *repository) GetById(id string) (*access_token.AccessToken, *errorUtils.RestErr) {
	session, err := cassandra.GetSession()

	if err != nil {
		return nil, errorUtils.NewInternalServerError(err.Error())
	}

	defer session.Close()

	var result access_token.AccessToken
	if err := session.Query(queryGetAccessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, errorUtils.NewNotFoundError("no access token found with given id")
		}
		return nil, errorUtils.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (r *repository) Create(at access_token.AccessToken) *errorUtils.RestErr {
	session, err := cassandra.GetSession()

	if err != nil {
		return errorUtils.NewInternalServerError(err.Error())
	}

	defer session.Close()

	if err := session.Query(queryCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).Exec(); err != nil {
		return errorUtils.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *repository) UpdateExpirationTime(at access_token.AccessToken) *errorUtils.RestErr {
	session, err := cassandra.GetSession()

	if err != nil {
		return errorUtils.NewInternalServerError(err.Error())
	}

	defer session.Close()

	if err := session.Query(queryUpdateExpires, at.Expires, at.AccessToken).Exec(); err != nil {
		return errorUtils.NewInternalServerError(err.Error())
	}

	return nil
}

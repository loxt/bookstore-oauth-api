package rest

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/loxt/bookstore-oauth-api/internal/domain/users"
	"github.com/loxt/bookstore-oauth-api/internal/utils/errors"
	"time"
)

var (
	usersRestClient = resty.New().SetBaseURL("http://localhost:8080").SetTimeout(100 * time.Millisecond).
		SetJSONMarshaler(json.Marshal).
		SetJSONUnmarshaler(json.Unmarshal)
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct{}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := usersRestClient.R().SetBody(request).Post("/users/login")

	if err != nil {
		return nil, errors.NewInternalServerError("invalid rest client response when trying to login user")
	}

	if response.StatusCode() > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Body(), &restErr)

		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}

		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users login response")
	}

	return &user, nil
}

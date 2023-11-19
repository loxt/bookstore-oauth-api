package rest

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases...")
	httpmock.ActivateNonDefault(usersRestClient.GetClient())
	defer httpmock.DeactivateAndReset()

	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	httpmock.RegisterMatcherResponder("POST", "http://localhost:8080/users/login",
		httpmock.BodyContainsString(`{"email": "email@gmail.com", "password": "password"}`),
		httpmock.NewStringResponder(-1, ""))

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, `{"message":"invalid login credentials", "status": "404", "error": "not_found"}`), nil
		})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, `{"message":"invalid login credentials", "status": 404, "error": "not_found"}`), nil
		})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"id": "1", "first_name": "Emannuel", "last_name": "M.", "email": "loxtdevanced@gmail.com"}`), nil
		})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"id": 1, "first_name": "Emannuel", "last_name": "M.", "email": "loxtdevanced@gmail.com"}`), nil
		})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Emannuel", user.FirstName)
	assert.EqualValues(t, "M.", user.LastName)
	assert.EqualValues(t, "loxtdevanced@gmail.com", user.Email)
}

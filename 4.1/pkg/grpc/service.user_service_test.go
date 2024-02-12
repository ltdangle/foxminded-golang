package grpc

import (
	"context"
	"grpc4_1/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s *UserService = NewUserService(model.NewUserRepoStub())

func TestCreateUser(t *testing.T) {
	// Assert error on empty request.
	_, err := s.CreateUser(context.Background(), &CreateUserRequest{})
	assert.NotNil(t, err)

	// Assert no error on valid request.
	_, err = s.CreateUser(context.Background(), &CreateUserRequest{
		Email:     "email@domain.net",
		Password:  "somepassword",
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Nickname:  "Nickname",
	})
	assert.Nil(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	// Assert error on empty request.
	_, err := s.GetUserByEmail(context.Background(), &GetUserByEmailRequest{})
	assert.NotNil(t, err)

	// Assert no error on valid request.
	_, err = s.GetUserByEmail(context.Background(), &GetUserByEmailRequest{
		Email: "email@domain.net",
	})
	assert.Nil(t, err)

}

func TestGetUserByID(t *testing.T) {
	// Assert error on empty request.
	_, err := s.GetUserByID(context.Background(), &GetUserByIDRequest{})
	assert.NotNil(t, err)

	// Assert no error on valid request.
	_, err = s.GetUserByID(context.Background(), &GetUserByIDRequest{
		Id: "someuuid",
	})
	assert.Nil(t, err)
}

func TestGetUsers(t *testing.T) {
	_, err := s.GetUsers(context.Background(), &GetUsersRequest{})
	assert.Nil(t, err)
}

func TestUpdateUser(t *testing.T) {
	// Assert error on empty request.
	_, err := s.UpdateUser(context.Background(), &UpdateUserRequest{})
	assert.NotNil(t, err)

	// Assert no error on valid request.
	_, err = s.UpdateUser(context.Background(), &UpdateUserRequest{
		Uuid:      "someuuid",
		Email:     "email@domain.net",
		Password:  "somepassword",
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Nickname:  "Nickname",
	})
	assert.Nil(t, err)
}

func TestDeleteUser(t *testing.T) {
	// Assert error on empty request.
	_, err := s.DeleteUser(context.Background(), &DeleteUserRequest{})
	assert.NotNil(t, err)

	// Assert no error on valid request.
	_, err = s.DeleteUser(context.Background(), &DeleteUserRequest{
		Id: "someuuid",
	})
	assert.Nil(t, err)
}

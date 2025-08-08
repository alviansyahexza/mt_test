package test

import (
	"errors"
	"testing"

	"github.com/alviansyahexza/mt_test/entity"
	"github.com/alviansyahexza/mt_test/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var userRepoMock = &UserRepoMock{Mock: mock.Mock{}}
var userService = service.NewUserService(userRepoMock)

func TestUserServiceSignUp(t *testing.T) {
	mockedUser := &entity.User{
		Id:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	userRepoMock.On("GetUserByEmail", mockedUser.Email).Return(nil, errors.New("user not found"))

	userRepoMock.On("CreateUser", mockedUser.Name, mockedUser.Email, mock.AnythingOfType("string")).Return(int64(mockedUser.Id), nil)

	userResult, err := userService.SignUp(mockedUser.Name, mockedUser.Email, mockedUser.Password)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, mockedUser.Id, userResult.Id)
	assert.Equal(t, mockedUser.Email, userResult.Email)
	assert.Equal(t, "", userResult.Password)

	userRepoMock.On("GetUserByEmail", mockedUser.Email).Unset()
	userRepoMock.On("GetUserByEmail", mockedUser.Email).Return(mockedUser, nil)
	userResult2, err2 := userService.SignUp(mockedUser.Name, mockedUser.Email, mockedUser.Password)
	assert.Nil(t, userResult2)
	assert.NotNil(t, err2)
	assert.Equal(t, "user already exists with this email", err2.Error())
}

func TestUserServiceSignIn(t *testing.T) {
	mockedPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(mockedPassword), bcrypt.DefaultCost)
	mockedUser := &entity.User{
		Id:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: string(hashedPassword),
	}

	userRepoMock.On("GetUserByEmail", mockedUser.Email).Return(mockedUser, nil)

	userResult, err := userService.SignIn(mockedUser.Email, mockedPassword)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, mockedUser.Id, userResult.Id)
	assert.Equal(t, mockedUser.Email, userResult.Email)
	assert.Equal(t, "", userResult.Password)

	userResultFailed, err := userService.SignIn(mockedUser.Email, "wrongpassword")
	assert.NotNil(t, err)
	assert.Nil(t, userResultFailed)
}

func TestUserServiceUpdateProfile(t *testing.T) {
	mockedUser := &entity.User{
		Id:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	userRepoMock.On("UpdateUser", mockedUser).Return(mockedUser, nil)

	_, err := userService.UpdateProfile(mockedUser)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	mockedUser.Name = ""
	failedUpdatedUser, err := userService.UpdateProfile(mockedUser)
	assert.Nil(t, failedUpdatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid data to update: name and id are required", err.Error())
	mockedUser.Id = 0
	failedUpdatedUser, err = userService.UpdateProfile(mockedUser)
	assert.Nil(t, failedUpdatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid data to update: name and id are required", err.Error())
}

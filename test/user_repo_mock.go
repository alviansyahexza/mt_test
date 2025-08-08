package test

import (
	"github.com/alviansyahexza/mt_test/entity"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) CreateUser(name, email, password string) (int64, error) {
	args := m.Called(name, email, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *UserRepoMock) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) UpdateUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if updatedUser, ok := args.Get(0).(*entity.User); ok {
		return updatedUser, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) GetUserById(id int) (*entity.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

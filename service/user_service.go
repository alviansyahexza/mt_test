package service

import (
	"errors"

	"github.com/alviansyahexza/mt_test/entity"
	"github.com/alviansyahexza/mt_test/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) SignUp(name, email, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	id, err := s.userRepo.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	user.Id = int(id)
	user.Password = ""
	return user, nil
}

func (s *UserService) SignIn(email, password string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) UpdateProfile(user *entity.User) (*entity.User, error) {
	if user.Name == "" || user.Id == 0 {
		return nil, errors.New("invalid data to update: name and id are required")
	}
	u, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

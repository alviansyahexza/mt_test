package service

import (
	"database/sql"

	"github.com/alviansyahexza/mt_test/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
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
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id "
	err = s.db.QueryRow(query, name, email, string(hashedPassword)).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) SignIn(email, password string) (*entity.User, error) {
	user := &entity.User{}
	query := "SELECT id, name, email, password FROM users WHERE email = $1 "
	err := s.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
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
	query := "UPDATE users SET name = $1 WHERE id = $2 Returning id "
	err := s.db.QueryRow(query, user.Name, user.Id).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	query = "SELECT id, name, email FROM users WHERE id = $1 "
	err = s.db.QueryRow(query, user.Id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

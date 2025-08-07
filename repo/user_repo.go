package repo

import (
	"github.com/alviansyahexza/mt_test/entity"
)

type UserRepo interface {
	CreateUser(name, email, password string) (int64, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
}

package repo

import (
	"database/sql"

	"github.com/alviansyahexza/mt_test/entity"
)

type userRepoImpl struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepoImpl {
	return &userRepoImpl{db: db}
}

func (r *userRepoImpl) CreateUser(name, email, password string) (int64, error) {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := r.db.QueryRow(query, name, email, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepoImpl) GetUserByEmail(email string) (*entity.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = $1"
	user := &entity.User{}
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepoImpl) UpdateUser(user *entity.User) (*entity.User, error) {
	query := "UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name, email"
	err := r.db.QueryRow(query, user.Name, user.Id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepoImpl) GetUserById(id int) (*entity.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE id = $1"
	user := &entity.User{}
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

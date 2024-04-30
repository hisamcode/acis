package repository

import "github.com/hisamcode/acis/internal/data"

type UserDatabaseRepo interface {
	Insert(user *data.User) error
	GetByEmail(email string) (*data.User, error)
	UpdateUser(user *data.User) error
}

package repository

import (
	"time"

	"github.com/hisamcode/acis/internal/data"
)

type UserDatabaseRepoer interface {
	Insert(user *data.User) error
	GetByEmail(email string) (*data.User, error)
	Update(user *data.User) error
	GetForToken(tokenScope, tokenPlaintext string) (*data.User, error)
}

type TokenDatabaseRepoer interface {
	New(userID int64, ttl time.Duration, scope string) (*data.Token, error)
	Insert(token *data.Token) error
	DeleteAllForUser(scope string, userID int64) error
}

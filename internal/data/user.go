package data

import (
	"errors"
	"time"

	"github.com/hisamcode/acis/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// represent individual user
type User struct {
	ID          int64
	CreatedAt   time.Time
	Name, Email string
	Password    password
	Activated   bool
	Version     int
}

// The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all, versus a plaintext password which is the empty string "".
type password struct {
	PlainText *string
	Hash      []byte
}

// generate bcrypt hash from plaintextPassword, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.PlainText = &plaintextPassword
	p.Hash = hash

	return nil
}

// checks whether the provided plaintext password matches the
// hashed password stored in the struct.
// returning true if it matches and false otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil

}

// check email != "" and check email from validator.EmailRX regexp
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "Email must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "Email must be valid email address")
}

// check password not empty and password >=8 and password <=72
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password must be provided")
	v.Check(len(password) >= 8, "password", "Password must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "Password must not be more than 72 bytes long")
}

// check user.Name not empty and check length user.name <=500 bytes
// and call ValidateEmail() and call ValidatePasswordPlaintext()
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "Name must be provided")
	v.Check(len(user.Name) <= 500, "name", "Name must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.PlainText != nil {
		ValidatePasswordPlaintext(v, *user.Password.PlainText)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}

}

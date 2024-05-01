package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/hisamcode/acis/internal/validator"
)

const (
	ScopeActivation = "activation"
)

// represent individual token
type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

// generate token without padding =,
// encode with base-32 string, which result 26 characters
// if we encoded the random bytes using hexadecimal (base-16) the string would be 32 characters long instead.
func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// fill the byte slice with random bytes from your operating system's CSPRING
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "token must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "token must be 26 bytes long")
}

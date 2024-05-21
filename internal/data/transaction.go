package data

import (
	"fmt"
	"time"

	"github.com/hisamcode/acis/internal/validator"
)

type Transaction struct {
	ID        int64
	CreatedAt time.Time
	Nominal   float64
	Detail    string
	WTSID     int8
	EmojiID   Emoji
	Version   int
	ProjectID int64
	CreatedBy int64
}

func ValidateTransaction(v *validator.Validator, transaction *Transaction) {
	limit := 1_000_000_000_000
	v.Check(int(transaction.Nominal) < limit, "nominal", fmt.Sprintf("nominal must be more than %d", limit))
}

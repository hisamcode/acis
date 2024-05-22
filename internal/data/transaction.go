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
	Emoji     Emoji
	Version   int
	ProjectID int64
	CreatedBy int64
}

func ValidateTransaction(v *validator.Validator, transaction *Transaction) {
	if len(transaction.Detail) > 0 {
		v.Check(len(transaction.Detail) < 500, "detail", "detail must not be more than 500 bytes/character long")
	}
	v.Check(transaction.Nominal != 0, "nominal", "nominal can't 0")
	if transaction.Nominal > float64(1_000_000) {
		limit := float64(1_000_000_000_000)
		v.Check(transaction.Nominal < limit, "nominal", fmt.Sprintf("nominal must be more than %f", limit))
	}
}

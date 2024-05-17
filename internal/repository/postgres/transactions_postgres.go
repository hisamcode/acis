package postgres

import (
	"context"
	"database/sql"

	"github.com/hisamcode/acis/internal/data"
)

type TransactionModel struct {
	DB *sql.DB
}

func (m TransactionModel) Insert(transaction *data.Transaction) error {
	query := `
	INSERT INTO transactions (nominal, detail, emoji_id, wts_id, project_id, created_by)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, version
	`

	args := []any{
		transaction.Nominal,
		transaction.Detail,
		transaction.EmojiID.Encode(),
		transaction.WTSID,
		transaction.ProjectID,
		transaction.CreatedBy,
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		return err.Err()
	}

	return nil
}

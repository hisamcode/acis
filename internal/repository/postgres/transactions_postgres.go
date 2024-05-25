package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/hisamcode/acis/internal/data"
)

type TransactionModel struct {
	DB *sql.DB
}

func (m TransactionModel) Insert(transaction *data.Transaction) error {
	query := `
	INSERT INTO transactions (nominal, detail, emoji_id, wts_id, project_id, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, version
	`

	args := []any{
		transaction.Nominal,
		transaction.Detail,
		transaction.Emoji.ID,
		transaction.WTSID,
		transaction.ProjectID,
		transaction.CreatedBy,
		transaction.CreatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (m TransactionModel) LatestBetweenDate(project data.Project, date1, date2 time.Time, limit int, lastDisplayedID int64) ([]data.Transaction, error) {
	var query string
	var args []any

	if lastDisplayedID < 1 {
		query = `
		SELECT 
			created_at, nominal, detail, emoji_id, wts_id, project_id 
		FROM transactions
		WHERE
			project_id = $1
			AND created_at >= $2 
			AND created_at < $3
		ORDER BY created_at DESC
		LIMIT $4
		`
		args = []any{project.ID, date1, date2, limit}
	} else {
		query = `
		SELECT 
			created_at, nominal, detail, emoji_id, wts_id, project_id 
		FROM transactions
		WHERE
			project_id = $1
			AND created_at >= $2 
			AND created_at < $3
			AND id < $4
		ORDER BY created_at DESC
		LIMIT $5
		`
		args = []any{project.ID, date1, date2, lastDisplayedID, limit}
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	transactions := []data.Transaction{}

	for rows.Next() {
		transaction := data.Transaction{}

		err := rows.Scan(
			&transaction.CreatedAt,
			&transaction.Nominal,
			&transaction.Detail,
			&transaction.Emoji.ID,
			&transaction.WTSID,
			&transaction.ProjectID,
		)

		if err != nil {
			return nil, err
		}

		emoji, err := project.FindEmoji(transaction.Emoji.ID)
		if err != nil {
			return nil, err
		}

		transaction.Emoji = emoji
		transactions = append(transactions, transaction)
	}

	// error during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

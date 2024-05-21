package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hisamcode/acis/internal/data"
	"github.com/lib/pq"
)

type ProjectModel struct {
	DB *sql.DB
}

// insert with return id and error
func (m ProjectModel) Insert(project *data.Project) (int64, error) {
	query := `
	INSERT INTO projects (name, detail, emojis, user_id)
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, version
	`

	emojisEncoded := []string{}
	for _, v := range project.Emojis {
		emojisEncoded = append(emojisEncoded, v.Encode())
	}

	args := []any{project.Name, project.Detail, pq.Array(emojisEncoded), project.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&project.ID,
		&project.CreatedAt,
		&project.Version,
	)
	if err != nil {
		return -1, err
	}

	return project.ID, nil
}

func (m ProjectModel) Get(id int64) (*data.Project, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	query := `
	SELECT id, name, detail, emojis, created_at, version, user_id
	FROM projects
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var project data.Project

	emojisEncoded := []string{}

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Detail,
		pq.Array(&emojisEncoded),
		&project.CreatedAt,
		&project.Version,
		&project.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for _, v := range emojisEncoded {
		emoji := data.Emoji{}
		err := emoji.Decode(v)
		if err != nil {
			return nil, err
		}
		project.Emojis = append(project.Emojis, emoji)
	}

	return &project, nil
}

func (m ProjectModel) LatestByUserID(userID int64) ([]data.Project, error) {

	query := `
	SELECT id, name, detail, emojis, created_at, version, user_id
	FROM projects
	WHERE user_id = $1
	ORDER BY id DESC LIMIT 10
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []data.Project{}

	for rows.Next() {
		var project data.Project
		emojis := []string{}

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Detail,
			pq.Array(&emojis),
			&project.CreatedAt,
			&project.Version,
			&project.UserID,
		)

		if err != nil {
			return nil, err
		}

		project.Emojis = data.DecodeEmojis(emojis)

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil

}

func (m ProjectModel) Update(project *data.Project) error {
	query := `
	UPDATE projects
	SET name = $1, detail = $2, emojis = $3, user_id = $4, version = version + 1
	WHERE id = $5 AND version = $6
	returning version
	`

	// encoded emoji
	encodedEmojis := []string{}
	for _, v := range project.Emojis {
		encodedEmojis = append(encodedEmojis, v.Encode())
	}

	args := []any{
		project.Name,
		project.Detail,
		pq.Array(encodedEmojis),
		project.UserID,
		project.ID,
		project.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Execute the sql query, if no matching row could be found, we know the project
	// version has changed (or the record has been deleted) and we return error
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&project.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

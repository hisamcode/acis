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
	if len(project.Categories) < 1 {
		project.Categories = append(project.Categories, "empty;empty")
	}
	query := `
	INSERT INTO projects (name, detail, categories, user_id)
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, version
	`
	args := []any{project.Name, project.Detail, pq.Array(project.Categories), project.UserID}

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
	SELECT id, name, detail, categories, created_at, version, user_id
	FROM projects
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var project data.Project

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Detail,
		pq.Array(&project.Categories),
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

	return &project, nil
}

func (m ProjectModel) LatestByUserID(userID int64) ([]data.Project, error) {

	query := `
	SELECT id, name, detail, categories, created_at, version, user_id
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

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Detail,
			pq.Array(&project.Categories),
			&project.CreatedAt,
			&project.Version,
			&project.UserID,
		)

		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil

}
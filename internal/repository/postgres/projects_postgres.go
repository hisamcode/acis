package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hisamcode/acis/internal/data"
)

type ProjectModel struct {
	DB *sql.DB
}

// insert with return id and error
func (m ProjectModel) Insert(project *data.Project) (int64, error) {
	query := `
	INSERT INTO projects (name, detail, user_id)
	VALUES ($1, $2, $3) 
	RETURNING id, created_at, version
	`
	args := []any{project.Name, project.Detail, project.UserID}

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

func (m ProjectModel) Get(id int) (*data.Project, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	query := `
	SELECT id, name, detail, created_at, version, user_id
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
	SELECT id, name, detail, created_at, version, user_id
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

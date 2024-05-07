package data

import (
	"time"

	"github.com/hisamcode/acis/internal/validator"
)

type Project struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Detail    string
	Version   int
	UserID    int64
}

func ValidateProject(v *validator.Validator, project *Project) {
	v.Check(project.Name != "", "name", "Name must be provided")
	v.Check(len(project.Name) <= 500, "name", "Name must not be more than 500 bytes/character long")
	v.Check(len(project.Detail) <= 2000, "detail", "Detail must not be more than 2000 bytes/character long")
}

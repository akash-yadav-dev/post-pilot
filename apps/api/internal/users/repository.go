package users

import (
	"context"
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) Create(ctx context.Context, email string, name string) (string, error) {

	query := `
	INSERT INTO users (email, name)
	VALUES ($1, $2)
	RETURNING id
	`

	var id string

	err := r.DB.QueryRowContext(ctx, query, email, name).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

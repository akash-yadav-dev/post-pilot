package users

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	Plan      string
	Credits   int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

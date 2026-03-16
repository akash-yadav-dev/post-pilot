package database

import "fmt"

type DB struct {
	Name     string
	Username string
	Password string
}

func NewDB(name string, username string, password string) (*DB, error) {
	if name == "" {
		return nil, fmt.Errorf("database name is required")
	}
	if username == "" {
		return nil, fmt.Errorf("database username is required")
	}

	return &DB{
		Name:     name,
		Username: username,
		Password: password,
	}, nil
}

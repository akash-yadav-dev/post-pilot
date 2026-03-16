package postgres

import (
	"context"
	"database/sql"
	"os"
	"time"
)

type DB struct {
	*sql.DB
}

type DBConfig struct {
	Connection      *Connection
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

func NewPostgresDB() (*DB, error) {
	conn := NewConnection(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "app_db"),
	)

	cfg := DBConfig{
		Connection:      conn,
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		PingTimeout:     5 * time.Second,
	}

	return NewWithConfig(cfg)
}

func NewWithConfig(cfg DBConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.Connection.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func getEnv(key, fallback string) string {

	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}

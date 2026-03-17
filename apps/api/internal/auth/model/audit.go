package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type AuditEvent struct {
	UserID       *uuid.UUID
	ActorType    string
	ActorIP      string
	UserAgent    string
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	Metadata     map[string]any
	Succeeded    bool
	ErrorMessage string
}

func NullableUUID(id *uuid.UUID) sql.NullString {
	if id == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: id.String(), Valid: true}
}

func NullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

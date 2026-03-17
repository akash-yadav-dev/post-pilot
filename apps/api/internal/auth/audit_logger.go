package auth

import (
	"context"
	"database/sql"
	"encoding/json"

	"post-pilot/apps/api/internal/auth/model"
	"post-pilot/packages/logger"
)

type AuditLogger interface {
	LogEvent(ctx context.Context, event model.AuditEvent)
}

type DBAuditLogger struct {
	db     *sql.DB
	logger logger.Logger
}

func NewDBAuditLogger(db *sql.DB, appLogger logger.Logger) *DBAuditLogger {
	return &DBAuditLogger{db: db, logger: appLogger}
}

func (l *DBAuditLogger) LogEvent(ctx context.Context, event model.AuditEvent) {
	metadataBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		metadataBytes = []byte(`{}`)
	}

	query := `
		INSERT INTO audit_logs (
			user_id,
			actor_type,
			actor_ip,
			user_agent,
			action,
			resource_type,
			resource_id,
			metadata,
			succeeded,
			error_message
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9,$10)
	`

	_, dbErr := l.db.ExecContext(
		ctx,
		query,
		model.NullableUUID(event.UserID),
		event.ActorType,
		model.NullableString(event.ActorIP),
		model.NullableString(event.UserAgent),
		event.Action,
		event.ResourceType,
		model.NullableUUID(event.ResourceID),
		string(metadataBytes),
		event.Succeeded,
		model.NullableString(event.ErrorMessage),
	)

	fields := []logger.Field{
		{Key: "event", Value: "auth_audit"},
		{Key: "action", Value: event.Action},
		{Key: "resource_type", Value: event.ResourceType},
		{Key: "actor_type", Value: event.ActorType},
		{Key: "succeeded", Value: event.Succeeded},
		{Key: "metadata", Value: event.Metadata},
	}
	if event.UserID != nil {
		fields = append(fields, logger.Field{Key: "user_id", Value: event.UserID.String()})
	}
	if event.ActorIP != "" {
		fields = append(fields, logger.Field{Key: "actor_ip", Value: event.ActorIP})
	}

	if dbErr != nil {
		l.logger.Error("failed to persist auth audit event", dbErr, fields...)
		return
	}

	l.logger.Info("auth audit event recorded", fields...)
}

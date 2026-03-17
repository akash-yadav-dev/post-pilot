package auth

import (
	"context"
	"errors"
	"testing"

	"post-pilot/apps/api/internal/auth/model"
	"post-pilot/packages/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type fakeLogger struct {
	infoCount  int
	errorCount int
}

func (f *fakeLogger) Debug(string, ...logger.Field) {}
func (f *fakeLogger) Info(string, ...logger.Field)  { f.infoCount++ }
func (f *fakeLogger) Warn(string, ...logger.Field)  {}
func (f *fakeLogger) Error(string, error, ...logger.Field) {
	f.errorCount++
}
func (f *fakeLogger) With(...logger.Field) logger.Logger { return f }
func (f *fakeLogger) Sync() error                        { return nil }

func TestDBAuditLoggerLogEventInsertSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	log := &fakeLogger{}
	audit := NewDBAuditLogger(db, log)

	uid := uuid.New()

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	audit.LogEvent(context.Background(), model.AuditEvent{
		UserID:       &uid,
		ActorType:    "user",
		ActorIP:      "127.0.0.1",
		UserAgent:    "test-agent",
		Action:       "auth.login",
		ResourceType: "session",
		ResourceID:   &uid,
		Metadata:     map[string]any{"email": "a@example.com"},
		Succeeded:    true,
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}

	if log.infoCount != 1 {
		t.Fatalf("infoCount = %d, want 1", log.infoCount)
	}
	if log.errorCount != 0 {
		t.Fatalf("errorCount = %d, want 0", log.errorCount)
	}
}

func TestDBAuditLoggerLogEventInsertFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	log := &fakeLogger{}
	audit := NewDBAuditLogger(db, log)

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("db down"))

	audit.LogEvent(context.Background(), model.AuditEvent{
		ActorType:    "anonymous",
		Action:       "auth.login",
		ResourceType: "session",
		Metadata:     map[string]any{"email": "a@example.com"},
		Succeeded:    false,
		ErrorMessage: "invalid credentials",
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}

	if log.errorCount != 1 {
		t.Fatalf("errorCount = %d, want 1", log.errorCount)
	}
}

package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"post-pilot/apps/api/internal/auth/model"
	"post-pilot/apps/api/internal/auth/service"
	"post-pilot/packages/security"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditLogger interface {
	LogEvent(ctx context.Context, event model.AuditEvent)
}

type AuthHandler struct {
	service *service.AuthService
	audit   AuditLogger
}

func NewAuthHandler(s *service.AuthService, audit AuditLogger) *AuthHandler {
	return &AuthHandler{service: s, audit: audit}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.register",
			ResourceType: "user",
			Metadata:     map[string]any{"email": strings.ToLower(strings.TrimSpace(req.Email))},
			Succeeded:    false,
			ErrorMessage: "invalid registration payload",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid registration payload"})
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		message := "registration failed"

		switch {
		case errors.Is(err, service.ErrEmailAlreadyRegistered):
			status = http.StatusConflict
			message = err.Error()
		}

		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.register",
			ResourceType: "user",
			Metadata:     map[string]any{"email": strings.ToLower(strings.TrimSpace(req.Email))},
			Succeeded:    false,
			ErrorMessage: message,
		})

		c.JSON(status, gin.H{"error": message})
		return
	}

	userID := resp.UserID
	h.auditEvent(c, model.AuditEvent{
		UserID:       &userID,
		ActorType:    "user",
		Action:       "auth.register",
		ResourceType: "user",
		ResourceID:   &userID,
		Metadata:     map[string]any{"email": resp.Email},
		Succeeded:    true,
	})

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.login",
			ResourceType: "session",
			Metadata:     map[string]any{"email": strings.ToLower(strings.TrimSpace(req.Email))},
			Succeeded:    false,
			ErrorMessage: "invalid login payload",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login payload"})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		message := "login failed"

		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			status = http.StatusUnauthorized
			message = err.Error()
		case errors.Is(err, service.ErrAccountLocked):
			status = http.StatusLocked
			message = err.Error()
		}

		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.login",
			ResourceType: "session",
			Metadata:     map[string]any{"email": strings.ToLower(strings.TrimSpace(req.Email))},
			Succeeded:    false,
			ErrorMessage: message,
		})

		c.JSON(status, gin.H{"error": message})
		return
	}

	userID := resp.UserID
	h.auditEvent(c, model.AuditEvent{
		UserID:       &userID,
		ActorType:    "user",
		Action:       "auth.login",
		ResourceType: "session",
		ResourceID:   &userID,
		Metadata:     map[string]any{"email": resp.Email},
		Succeeded:    true,
	})

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req model.GoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.google_login",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: "invalid google login payload",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid google login payload"})
		return
	}

	resp, err := h.service.LoginWithGoogle(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		message := "google login failed"

		switch {
		case errors.Is(err, service.ErrInvalidGoogleToken):
			status = http.StatusUnauthorized
			message = err.Error()
		case errors.Is(err, service.ErrGoogleAuthDisabled):
			status = http.StatusNotImplemented
			message = err.Error()
		}

		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.google_login",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: message,
		})

		c.JSON(status, gin.H{"error": message})
		return
	}

	userID := resp.UserID
	h.auditEvent(c, model.AuditEvent{
		UserID:       &userID,
		ActorType:    "user",
		Action:       "auth.google_login",
		ResourceType: "session",
		ResourceID:   &userID,
		Metadata:     map[string]any{"email": resp.Email, "provider": "google"},
		Succeeded:    true,
	})

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.refresh",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: "invalid refresh payload",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh payload"})
		return
	}

	resp, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		message := "invalid refresh token"

		switch {
		case errors.Is(err, security.ErrExpiredToken):
			message = err.Error()
		case errors.Is(err, security.ErrTokenRevoked):
			message = err.Error()
		case errors.Is(err, security.ErrInvalidToken):
			message = err.Error()
		case errors.Is(err, service.ErrUserNotFound):
			status = http.StatusNotFound
			message = err.Error()
		}

		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.refresh",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: message,
		})

		c.JSON(status, gin.H{"error": message})
		return
	}

	userID := resp.UserID
	h.auditEvent(c, model.AuditEvent{
		UserID:       &userID,
		ActorType:    "user",
		Action:       "auth.refresh",
		ResourceType: "session",
		ResourceID:   &userID,
		Metadata:     map[string]any{"email": resp.Email},
		Succeeded:    true,
	})

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req model.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.logout",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: "invalid logout payload",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid logout payload"})
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.logout",
			ResourceType: "session",
			Succeeded:    false,
			ErrorMessage: "invalid refresh token",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	h.auditEvent(c, model.AuditEvent{
		ActorType:    "user",
		Action:       "auth.logout",
		ResourceType: "session",
		Succeeded:    true,
	})

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := strings.TrimSpace(c.GetString("auth_user_id"))
	if userID == "" {
		h.auditEvent(c, model.AuditEvent{
			ActorType:    "anonymous",
			Action:       "auth.me",
			ResourceType: "user",
			Succeeded:    false,
			ErrorMessage: "unauthorized",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to fetch profile"

		switch {
		case errors.Is(err, service.ErrUserNotFound):
			status = http.StatusNotFound
			message = err.Error()
		}

		h.auditEvent(c, model.AuditEvent{
			ActorType:    "user",
			Action:       "auth.me",
			ResourceType: "user",
			Succeeded:    false,
			ErrorMessage: message,
			Metadata:     map[string]any{"auth_user_id": userID},
		})

		c.JSON(status, gin.H{"error": message})
		return
	}

	uid := resp.UserID
	h.auditEvent(c, model.AuditEvent{
		UserID:       &uid,
		ActorType:    "user",
		Action:       "auth.me",
		ResourceType: "user",
		ResourceID:   &uid,
		Succeeded:    true,
	})

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) auditEvent(c *gin.Context, event model.AuditEvent) {
	if h.audit == nil {
		return
	}

	event.ActorIP = c.ClientIP()
	event.UserAgent = strings.TrimSpace(c.Request.UserAgent())

	if event.Metadata == nil {
		event.Metadata = map[string]any{}
	}
	event.Metadata["path"] = c.FullPath()
	event.Metadata["method"] = c.Request.Method

	if event.UserID == nil {
		if id := strings.TrimSpace(c.GetString("auth_user_id")); id != "" {
			if parsed, err := uuid.Parse(id); err == nil {
				event.UserID = &parsed
			}
		}
	}

	h.audit.LogEvent(c.Request.Context(), event)
}

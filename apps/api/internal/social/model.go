package social

import (
	"time"

	"github.com/google/uuid"
)

type SocialAccount struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	Platform       string     `json:"platform"`
	AccountID      string     `json:"account_id"`
	AccountName    string     `json:"account_name,omitempty"`
	AccountURL     string     `json:"account_url,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	TokenScope     string     `json:"token_scope,omitempty"`
	Status         string     `json:"status"`
	LastUsedAt     *time.Time `json:"last_used_at,omitempty"`
	LastError      string     `json:"last_error,omitempty"`
	ErrorCount     int        `json:"error_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ConnectAccountRequest struct {
	Platform       string     `json:"platform" binding:"required,oneof=twitter linkedin mastodon bluesky"`
	AccountID      string     `json:"account_id" binding:"required,min=1,max=255"`
	AccountName    string     `json:"account_name" binding:"omitempty,max=255"`
	AccountURL     string     `json:"account_url" binding:"omitempty,url"`
	AccessToken    string     `json:"access_token" binding:"required,min=1"`
	RefreshToken   string     `json:"refresh_token" binding:"omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at" binding:"omitempty"`
	TokenScope     string     `json:"token_scope" binding:"omitempty,max=1024"`
	Metadata       string     `json:"metadata" binding:"omitempty"`
}

type PublishPostResponse struct {
	PostID  uuid.UUID               `json:"post_id"`
	Status  string                  `json:"status"`
	Results []PublishTargetResponse `json:"results"`
	Errors  int                     `json:"errors"`
	Success int                     `json:"success"`
}

type PublishTargetResponse struct {
	TargetID        uuid.UUID `json:"target_id"`
	Platform        string    `json:"platform"`
	SocialAccountID uuid.UUID `json:"social_account_id"`
	Status          string    `json:"status"`
	ExternalPostID  string    `json:"external_post_id,omitempty"`
	ExternalPostURL string    `json:"external_post_url,omitempty"`
	Error           string    `json:"error,omitempty"`
}

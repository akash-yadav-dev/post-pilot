package model

import "github.com/google/uuid"

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthResponse struct {
	UserID uuid.UUID  `json:"user_id"`
	Name   string     `json:"name"`
	Email  string     `json:"email"`
	Tokens AuthTokens `json:"tokens"`
}

type ProfileResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
}

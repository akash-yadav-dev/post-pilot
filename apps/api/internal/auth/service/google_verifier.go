package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/idtoken"
)

type GoogleIdentity struct {
	Subject       string
	Email         string
	Name          string
	EmailVerified bool
}

type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleIdentity, error)
}

type GoogleIDTokenVerifier struct {
	clientID string
}

func NewGoogleIDTokenVerifier(clientID string) *GoogleIDTokenVerifier {
	return &GoogleIDTokenVerifier{clientID: strings.TrimSpace(clientID)}
}

func (v *GoogleIDTokenVerifier) Verify(ctx context.Context, idToken string) (*GoogleIdentity, error) {
	if strings.TrimSpace(v.clientID) == "" {
		return nil, errors.New("google oauth client id not configured")
	}

	payload, err := idtoken.Validate(ctx, strings.TrimSpace(idToken), v.clientID)
	if err != nil {
		return nil, fmt.Errorf("validate google id token: %w", err)
	}

	subject, _ := payload.Claims["sub"].(string)
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	emailVerified := false
	switch value := payload.Claims["email_verified"].(type) {
	case bool:
		emailVerified = value
	case string:
		emailVerified = strings.EqualFold(value, "true")
	}

	if strings.TrimSpace(subject) == "" || strings.TrimSpace(email) == "" {
		return nil, errors.New("google token missing required claims")
	}

	if strings.TrimSpace(name) == "" {
		name = strings.Split(strings.TrimSpace(email), "@")[0]
	}

	return &GoogleIdentity{
		Subject:       strings.TrimSpace(subject),
		Email:         strings.TrimSpace(email),
		Name:          strings.TrimSpace(name),
		EmailVerified: emailVerified,
	}, nil
}

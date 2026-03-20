package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"post-pilot/apps/api/internal/auth/model"
	"post-pilot/apps/api/internal/auth/repository"
	"post-pilot/packages/security"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrAccountLocked          = errors.New("account is temporarily locked")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidGoogleToken     = errors.New("invalid google token")
	ErrGoogleAuthDisabled     = errors.New("google auth is not configured")
)

type AuthService struct {
	repo           AuthRepository
	passwordSvc    *security.PasswordService
	jwtSvc         *security.JWTService
	googleVerifier GoogleTokenVerifier
	accessTokenTTL time.Duration
}

func NewAuthService(
	repo AuthRepository,
	passwordSvc *security.PasswordService,
	jwtSvc *security.JWTService,
	googleVerifier GoogleTokenVerifier,
	accessTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		repo:           repo,
		passwordSvc:    passwordSvc,
		jwtSvc:         jwtSvc,
		googleVerifier: googleVerifier,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)

	passwordHash, err := s.passwordSvc.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.CreateUserWithPassword(ctx, name, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	identity, err := s.repo.GetPasswordIdentityByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if identity.LockedUntil != nil && identity.LockedUntil.After(time.Now()) {
		return nil, ErrAccountLocked
	}

	validPassword, err := s.passwordSvc.VerifyPassword(req.Password, identity.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}

	if !validPassword {
		_ = s.repo.RecordFailedLogin(ctx, identity.ID)
		return nil, ErrInvalidCredentials
	}

	if err := s.repo.ResetFailedLogin(ctx, identity.ID); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateLastLogin(ctx, identity.ID); err != nil {
		return nil, err
	}

	if s.passwordSvc.NeedsRehash(identity.PasswordHash) {
		newHash, hashErr := s.passwordSvc.HashPassword(req.Password)
		if hashErr == nil {
			_ = s.repo.UpdatePasswordHash(ctx, identity.ID, newHash)
		}
	}

	user := &model.User{
		ID:    identity.ID,
		Name:  identity.Name,
		Email: identity.Email,
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (*model.AuthResponse, error) {
	if s.googleVerifier == nil {
		return nil, ErrGoogleAuthDisabled
	}

	identity, err := s.googleVerifier.Verify(ctx, req.IDToken)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not configured") {
			return nil, ErrGoogleAuthDisabled
		}
		return nil, ErrInvalidGoogleToken
	}

	user, err := s.repo.FindUserByProviderIdentity(ctx, "google", identity.Subject)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}

		user, err = s.repo.CreateOrLinkGoogleUser(
			ctx,
			identity.Name,
			normalizeEmail(identity.Email),
			identity.Subject,
			identity.EmailVerified,
		)
		if err != nil {
			if errors.Is(err, repository.ErrEmailAlreadyExists) {
				return nil, ErrEmailAlreadyRegistered
			}
			return nil, err
		}
	}

	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	rotatedRefreshToken, claims, err := s.jwtSvc.RotateRefreshToken(ctx, strings.TrimSpace(refreshToken))
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id from refresh token: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	accessToken, err := s.jwtSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Tokens: model.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: rotatedRefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(s.accessTokenTTL.Seconds()),
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.jwtSvc.RevokeRefreshToken(ctx, strings.TrimSpace(refreshToken))
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.ProfileResponse, error) {
	parsedID, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &model.ProfileResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (s *AuthService) buildAuthResponse(ctx context.Context, user *model.User) (*model.AuthResponse, error) {
	accessToken, err := s.jwtSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Tokens: model.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(s.accessTokenTTL.Seconds()),
		},
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

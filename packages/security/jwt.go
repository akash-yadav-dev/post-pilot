package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Sentinel errors for typed error handling at the call site.
var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrWrongTokenType = errors.New("wrong token type")
	ErrTokenRevoked   = errors.New("token has been revoked")
	ErrExpiredToken   = errors.New("token has expired")
)

// TokenType distinguishes access tokens from refresh tokens inside the claims,
// preventing a refresh token from being accepted where an access token is expected.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenStore is the interface your Redis (or any store) adapter must satisfy.
// It is used to support refresh-token rotation and revocation via jti tracking.
type TokenStore interface {
	// SaveRefreshToken persists a jti so it can later be validated or revoked.
	SaveRefreshToken(ctx context.Context, jti string, ttl time.Duration) error
	// ExistsRefreshToken returns true if the jti is still valid (not rotated/revoked).
	ExistsRefreshToken(ctx context.Context, jti string) (bool, error)
	// DeleteRefreshToken invalidates a jti (called on rotation or logout).
	DeleteRefreshToken(ctx context.Context, jti string) error
}

// JWTConfig holds all configuration needed to build a JWTService.
// Prefer loading these values from a secrets manager (e.g. AWS Secrets Manager,
// HashiCorp Vault) rather than plain environment variables.
type JWTConfig struct {
	// AccessSecret and RefreshSecret must be independent, cryptographically
	// random byte slices of at least 32 bytes each.
	AccessSecret  string
	RefreshSecret string

	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// JWTService issues and validates JWTs for the application.
// Access and refresh tokens use separate secrets so that a leaked refresh secret
// does not compromise access token validation (and vice-versa).
type JWTService struct {
	accessSecret           []byte
	refreshSecret          []byte
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
	store                  TokenStore
}

// JWTClaims extends the standard registered claims with application-specific fields.
type JWTClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// NewJWTService constructs a JWTService. store may be nil only in environments
// where refresh-token revocation is not required (e.g. unit tests); passing nil
// in production disables revocation checks and rotation.
func NewJWTService(cfg JWTConfig, store TokenStore) (*JWTService, error) {
	if len(cfg.AccessSecret) < 32 {
		return nil, errors.New("access secret must be at least 32 bytes")
	}
	if len(cfg.RefreshSecret) < 32 {
		return nil, errors.New("refresh secret must be at least 32 bytes")
	}
	if cfg.AccessSecret == cfg.RefreshSecret {
		return nil, errors.New("access and refresh secrets must be different")
	}

	return &JWTService{
		accessSecret:           []byte(cfg.AccessSecret),
		refreshSecret:          []byte(cfg.RefreshSecret),
		accessTokenExpiration:  cfg.AccessTokenExpiration,
		refreshTokenExpiration: cfg.RefreshTokenExpiration,
		store:                  store,
	}, nil
}

// GenerateAccessToken mints a short-lived access token for the given user.
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return s.generateToken(userID, email, AccessToken, s.accessSecret, s.accessTokenExpiration)
}

// GenerateRefreshToken mints a long-lived refresh token and, when a TokenStore
// is configured, persists the jti so it can later be rotated or revoked.
func (s *JWTService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID, email string) (string, error) {
	jti := uuid.NewString()

	tokenString, err := s.generateTokenWithJTI(jti, userID, email, RefreshToken, s.refreshSecret, s.refreshTokenExpiration)
	if err != nil {
		return "", err
	}

	if s.store != nil {
		if err := s.store.SaveRefreshToken(ctx, jti, s.refreshTokenExpiration); err != nil {
			return "", fmt.Errorf("persisting refresh token jti: %w", err)
		}
	}

	return tokenString, nil
}

// ValidateAccessToken parses and validates an access token, enforcing:
//   - correct signing algorithm (HS256)
//   - correct audience ("post-pilot-api")
//   - token_type == "access"
//   - expiry
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := s.parseToken(tokenString, s.accessSecret, "post-pilot-api")
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, ErrWrongTokenType
	}

	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token, enforcing:
//   - correct signing algorithm (HS256)
//   - correct audience ("post-pilot-refresh")
//   - token_type == "refresh"
//   - expiry
//   - jti exists in the store (not yet rotated or revoked)
func (s *JWTService) ValidateRefreshToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	claims, err := s.parseToken(tokenString, s.refreshSecret, "post-pilot-refresh")
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, ErrWrongTokenType
	}

	if s.store != nil {
		exists, err := s.store.ExistsRefreshToken(ctx, claims.ID)
		if err != nil {
			return nil, fmt.Errorf("checking refresh token revocation: %w", err)
		}
		if !exists {
			return nil, ErrTokenRevoked
		}
	}

	return claims, nil
}

// RotateRefreshToken validates the old refresh token, revokes it, and issues a
// new one — implementing single-use refresh token rotation.
// If the old token has already been used (jti missing from store), it returns
// ErrTokenRevoked, which should trigger a full logout on the client.
func (s *JWTService) RotateRefreshToken(ctx context.Context, oldTokenString string) (string, *JWTClaims, error) {
	claims, err := s.ValidateRefreshToken(ctx, oldTokenString)
	if err != nil {
		return "", nil, err
	}

	// Revoke the old jti before issuing the new token.
	if s.store != nil {
		if err := s.store.DeleteRefreshToken(ctx, claims.ID); err != nil {
			return "", nil, fmt.Errorf("revoking old refresh token: %w", err)
		}
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", nil, fmt.Errorf("parsing user id from claims: %w", err)
	}

	newToken, err := s.GenerateRefreshToken(ctx, userID, claims.Email)
	if err != nil {
		return "", nil, err
	}

	return newToken, claims, nil
}

// RevokeRefreshToken explicitly invalidates a refresh token (e.g. on logout).
func (s *JWTService) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	// Parse without revocation check — we want to revoke even if it's somehow
	// already absent from the store.
	claims, err := s.parseToken(tokenString, s.refreshSecret, "post-pilot-refresh")
	if err != nil {
		return err
	}

	if s.store == nil {
		return nil
	}

	return s.store.DeleteRefreshToken(ctx, claims.ID)
}

// RefreshTokenExpiration exposes the configured TTL (e.g. for setting cookie Max-Age).
func (s *JWTService) RefreshTokenExpiration() time.Duration {
	return s.refreshTokenExpiration
}

// --------------------------------------------------------------------------
// Internal helpers
// --------------------------------------------------------------------------

func (s *JWTService) generateToken(
	userID uuid.UUID,
	email string,
	tokenType TokenType,
	secret []byte,
	expiration time.Duration,
) (string, error) {
	return s.generateTokenWithJTI(uuid.NewString(), userID, email, tokenType, secret, expiration)
}

func (s *JWTService) generateTokenWithJTI(
	jti string,
	userID uuid.UUID,
	email string,
	tokenType TokenType,
	secret []byte,
	expiration time.Duration,
) (string, error) {
	now := time.Now()

	audience := "post-pilot-api"
	if tokenType == RefreshToken {
		audience = "post-pilot-refresh"
	}

	claims := JWTClaims{
		UserID:    userID.String(),
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID.String(),
			Issuer:    "post-pilot",
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func (s *JWTService) parseToken(tokenString string, secret []byte, expectedAudience string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Algorithm pinning — reject anything that isn't HS256.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return secret, nil
		},
		jwt.WithAudience(expectedAudience), // enforce aud claim
		jwt.WithExpirationRequired(),       // reject tokens without exp
		jwt.WithIssuedAt(),                 // reject tokens with future iat
	)

	if err != nil {
		// Surface expiry as a distinct error so callers can react appropriately
		// (e.g. redirect to /refresh vs. redirect to /login).
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %s", ErrInvalidToken, err.Error())
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

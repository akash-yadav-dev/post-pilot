package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Sentinel errors allow callers to handle specific failure cases
// without string matching.
var (
	ErrEmptyPassword   = errors.New("password cannot be empty")
	ErrPasswordTooLong = errors.New("password exceeds maximum length")
	ErrInvalidHash     = errors.New("password hash is invalid or corrupt")
)

const (
	// MinCost is the minimum acceptable bcrypt work factor for production.
	// OWASP recommends >= 12 as of 2024.
	MinCost = 12

	// maxPasswordBytes is bcrypt's hard truncation limit.
	// Passwords longer than this silently hash to the same value as
	// their first 72 bytes, which is a correctness bug.
	maxPasswordBytes = 72
)

// PasswordService hashes and verifies passwords using bcrypt.
type PasswordService struct {
	cost int
}

// NewPasswordService constructs a PasswordService with the given bcrypt cost.
// Returns an error if cost is outside the safe operating range.
func NewPasswordService(cost int) (*PasswordService, error) {
	if cost < MinCost {
		return nil, fmt.Errorf(
			"bcrypt cost %d is below the minimum of %d; raise it to meet OWASP recommendations",
			cost, MinCost,
		)
	}
	if cost > bcrypt.MaxCost {
		return nil, fmt.Errorf(
			"bcrypt cost %d exceeds bcrypt.MaxCost (%d)",
			cost, bcrypt.MaxCost,
		)
	}
	return &PasswordService{cost: cost}, nil
}

// HashPassword returns a bcrypt hash of password.
//
// Rejects empty passwords and passwords longer than 72 bytes — bcrypt
// silently truncates at that boundary, meaning a 200-char password and
// its first 72 chars would produce the same hash.
func (s *PasswordService) HashPassword(password string) (string, error) {
	if err := validatePassword(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword checks whether password matches the stored hash.
//
// Returns:
//   - (true,  nil)   — password is correct
//   - (false, nil)   — password is wrong (not an error; a result)
//   - (false, error) — hash is malformed or bcrypt failed unexpectedly
//
// Callers should distinguish the error case from the mismatch case to
// avoid silently ignoring hash corruption.
func (s *PasswordService) VerifyPassword(password, hash string) (bool, error) {
	if err := validatePassword(password); err != nil {
		return false, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		// Wrong password — this is a normal login-failure result, not an error.
		return false, nil
	}

	// Any other error (malformed hash, unsupported version, etc.) is a
	// real problem that should be logged and alerted on.
	return false, fmt.Errorf("%w: %s", ErrInvalidHash, err.Error())
}

// NeedsRehash reports whether the stored hash was produced at a lower cost
// than the service is currently configured for, or is otherwise unparseable.
//
// Call this after a successful VerifyPassword. If it returns true, hash the
// plaintext password again with HashPassword and persist the updated hash.
// This ensures all active users are silently upgraded whenever you raise the
// cost factor.
func (s *PasswordService) NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		// Malformed hash — force a rehash so the record gets fixed on next login.
		return true
	}
	return cost < s.cost
}

// Cost returns the bcrypt work factor this service is configured to use.
// Useful for logging and observability.
func (s *PasswordService) Cost() int {
	return s.cost
}

// validatePassword enforces the two invariants that bcrypt cannot check itself.
func validatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}
	// len([]byte(password)) counts bytes, not runes — correct for bcrypt's limit.
	if len([]byte(password)) > maxPasswordBytes {
		return fmt.Errorf("%w: bcrypt truncates at %d bytes", ErrPasswordTooLong, maxPasswordBytes)
	}
	return nil
}

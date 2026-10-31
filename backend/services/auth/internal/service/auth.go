// Package service implements the auth business rules: registration, login,
// token validation/revocation and password reset.
package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/repository"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// specialChars mirrors the frontend password requirements.
const specialChars = "!@#$%^&*()_+-={}[]:;<>,.?/~|"

// AuthService combines the account repositories.
type AuthService struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
	tokens   *auth.TokenManager
	now      func() time.Time
	resetTTL time.Duration
}

// NewAuthService wires the auth business rules.
func NewAuthService(users repository.UserRepository, sessions repository.SessionRepository, tokens *auth.TokenManager, resetTTL time.Duration) *AuthService {
	return &AuthService{users: users, sessions: sessions, tokens: tokens, now: time.Now, resetTTL: resetTTL}
}

// Register validates and creates a new account, then returns a session.
func (s *AuthService) Register(ctx context.Context, input domain.RegisterInput) (domain.User, string, error) {
	input.Sanitize()
	if err := validateEmail(input.Email); err != nil {
		return domain.User{}, "", err
	}
	if input.Name == "" {
		input.Name = localPart(input.Email)
	}
	if err := validateName(input.Name); err != nil {
		return domain.User{}, "", err
	}
	if err := validatePassword(input.Password); err != nil {
		return domain.User{}, "", err
	}
	if input.Role == "" {
		input.Role = domain.RoleCustomer
	}
	if !input.Role.Valid() {
		return domain.User{}, "", domain.Validationf("role must be one of: admin, staff, customer")
	}
	if input.Role != domain.RoleCustomer {
		// Self-registration only creates customer accounts; other roles are
		// provisioned by an admin through the management flow.
		return domain.User{}, "", domain.Forbiddenf("self-registration is limited to customer accounts")
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return domain.User{}, "", err
	}
	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         input.Role,
		Status:       domain.UserActive,
		CompanyName:  input.CompanyName,
	}
	if err := s.users.Create(ctx, &user); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return domain.User{}, "", domain.Conflictf("an account with this email address already exists. Please sign in instead")
		}
		return domain.User{}, "", err
	}
	token, err := s.issueSession(ctx, user, "", "")
	if err != nil {
		return domain.User{}, "", err
	}
	return user, token, nil
}

// Login verifies credentials and returns a session.
func (s *AuthService) Login(ctx context.Context, input domain.LoginInput) (domain.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := input.Password
	if email == "" || password == "" {
		return domain.User{}, "", domain.Unauthorized("Invalid email or password. Please verify your credentials.")
	}
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, "", domain.Unauthorized("Invalid email or password. Please verify your credentials.")
		}
		return domain.User{}, "", err
	}
	if user.Status != domain.UserActive || !auth.VerifyPassword(user.PasswordHash, password) {
		return domain.User{}, "", domain.Unauthorized("Invalid email or password. Please verify your credentials.")
	}
	now := s.now().UTC()
	if err := s.users.TouchLogin(ctx, user.ID, now); err != nil {
		return domain.User{}, "", err
	}
	token, err := s.issueSession(ctx, *user, "", "")
	if err != nil {
		return domain.User{}, "", err
	}
	return *user, token, nil
}

// Validate verifies a bearer token and returns the active user. The session is
// checked so a revoked (logged-out) token is rejected.
func (s *AuthService) Validate(ctx context.Context, token string) (domain.User, error) {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	user, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	if user.Status != domain.UserActive {
		return domain.User{}, domain.ErrUnauthorized
	}
	if claims.SessionID != "" {
		session, err := s.sessions.GetSessionByID(ctx, claims.SessionID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.User{}, domain.ErrUnauthorized
			}
			return domain.User{}, err
		}
		if session.RevokedAt != nil || !session.ExpiresAt.After(s.now()) {
			return domain.User{}, domain.ErrUnauthorized
		}
	}
	return *user, nil
}

// Logout revokes the active session so the bearer token stops validating.
func (s *AuthService) Logout(ctx context.Context, token string) error {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		// Invalid/expired tokens are already unusable; treat logout as success.
		return nil
	}
	if claims.SessionID != "" {
		if err := s.sessions.RevokeSession(ctx, claims.SessionID, s.now()); err != nil && !errors.Is(err, domain.ErrNotFound) {
			return err
		}
	}
	return nil
}

// ResetPassword issues a password-reset token for the user if one exists. It
// always returns success to avoid leaking which emails have accounts.
func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if err := validateEmail(email); err != nil {
		return err
	}
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	raw, err := auth.RandomToken(32)
	if err != nil {
		return err
	}
	token := &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(raw),
		ExpiresAt: s.now().Add(s.resetTTL),
	}
	return s.users.CreatePasswordResetToken(ctx, token)
}

// issueSession creates a session row, embeds its id in the access token and
// returns the signed token.
func (s *AuthService) issueSession(ctx context.Context, user domain.User, userAgent, ip string) (string, error) {
	refresh, err := auth.RandomToken(32)
	if err != nil {
		return "", err
	}
	now := s.now().UTC()
	session := &domain.Session{
		UserID:           user.ID,
		RefreshTokenHash: auth.HashToken(refresh),
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        now.Add(30 * 24 * time.Hour),
	}
	if err := s.sessions.CreateSession(ctx, session); err != nil {
		return "", err
	}
	return s.tokens.Issue(user, session.ID)
}

// --- validation helpers -----------------------------------------------------

func validateEmail(email string) error {
	if email == "" {
		return domain.Validationf("email is required")
	}
	if !emailPattern.MatchString(email) {
		return domain.Validationf("please enter a valid email address")
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 2 {
		return domain.Validationf("name must be at least 2 characters long")
	}
	if len(name) > 100 {
		return domain.Validationf("name must not exceed 100 characters")
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return domain.Validationf("password is required")
	}
	if len(password) <= 8 {
		return domain.Validationf("password must have more than 8 characters")
	}
	if !hasUpper(password) {
		return domain.Validationf("password must contain at least one uppercase letter")
	}
	if !hasLower(password) {
		return domain.Validationf("password must contain at least one lowercase letter")
	}
	if !hasSpecial(password) {
		return domain.Validationf("password must contain at least one special character")
	}
	return nil
}

func hasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLower(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasSpecial(s string) bool {
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

func localPart(email string) string {
	at := strings.Index(email, "@")
	if at < 0 {
		return email
	}
	return email[:at]
}

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/repository"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/validation"
)

// AuthService contains account registration and login rules.
type AuthService struct {
	users  repository.UserRepository
	tokens *auth.TokenManager
}

// NewAuthService wires the auth service to its dependencies.
func NewAuthService(users repository.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// RegisterInput is the data required to create a customer account.
type RegisterInput struct {
	Email    string
	Password string
	FullName string
}

// Register creates a new customer account. Staff and admin accounts are
// created separately by an admin.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	email := normalizeEmail(input.Email)
	if err := validation.Email(email); err != nil {
		return nil, err
	}
	if err := validation.Password(input.Password); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		FullName:     strings.TrimSpace(input.FullName),
		Role:         domain.RoleCustomer,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies credentials and returns an access token plus the user.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.users.ByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", nil, domain.ErrUnauthorized
		}
		return "", nil, err
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		return "", nil, domain.ErrUnauthorized
	}

	token, err := s.IssueToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// IssueToken creates an access token for an existing user.
func (s *AuthService) IssueToken(user *domain.User) (string, error) {
	token, err := s.tokens.Generate(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

// ByID loads a user by ID, used by the authenticated profile endpoint.
func (s *AuthService) ByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.ByID(ctx, id)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

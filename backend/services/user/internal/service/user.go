// Package service implements the user business rules: profile management,
// staff administration, preferences and account deletion.
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/repository"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// UserService implements profile and staff-administration rules.
type UserService struct {
	users repository.UserRepository
	now   func() time.Time
}

// NewUserService wires the user business rules.
func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users, now: time.Now}
}

// Me returns the current user's profile.
func (s *UserService) Me(ctx context.Context, actor domain.Principal) (*domain.User, error) {
	if actor.UserID == "" {
		return nil, domain.ErrUnauthorized
	}
	return s.users.GetByID(ctx, actor.UserID)
}

// UpdateProfile applies self-service profile changes.
func (s *UserService) UpdateProfile(ctx context.Context, actor domain.Principal, patch domain.UpdateProfilePatch) (*domain.User, error) {
	if actor.UserID == "" {
		return nil, domain.ErrUnauthorized
	}
	current, err := s.users.GetByID(ctx, actor.UserID)
	if err != nil {
		return nil, err
	}
	patch.Name = strings.TrimSpace(patch.Name)
	if patch.Name != "" {
		if len(patch.Name) < 2 {
			return nil, domain.Validationf("name must be at least 2 characters long")
		}
		current.Name = patch.Name
	}
	current.CompanyName = strings.TrimSpace(patch.CompanyName)
	current.Phone = strings.TrimSpace(patch.Phone)
	current.AvatarURL = strings.TrimSpace(patch.AvatarURL)
	return s.users.UpdateProfile(ctx, current)
}

// Preferences returns the current user's preferences.
func (s *UserService) Preferences(ctx context.Context, actor domain.Principal) (map[string]any, error) {
	user, err := s.Me(ctx, actor)
	if err != nil {
		return nil, err
	}
	if user.Preferences == nil {
		return map[string]any{}, nil
	}
	return user.Preferences, nil
}

// UpdatePreferences merges preference updates over the existing set.
func (s *UserService) UpdatePreferences(ctx context.Context, actor domain.Principal, updates map[string]any) (map[string]any, error) {
	if actor.UserID == "" {
		return nil, domain.ErrUnauthorized
	}
	current, err := s.Preferences(ctx, actor)
	if err != nil {
		return nil, err
	}
	for key, value := range updates {
		current[key] = value
	}
	if err := s.users.UpdatePreferences(ctx, actor.UserID, current); err != nil {
		return nil, err
	}
	return current, nil
}

// DeleteAccount deactivates the current user's account.
func (s *UserService) DeleteAccount(ctx context.Context, actor domain.Principal) error {
	if actor.UserID == "" {
		return domain.ErrUnauthorized
	}
	if err := s.users.SetStatus(ctx, actor.UserID, domain.UserInactive, true); err != nil {
		return err
	}
	return nil
}

// ListStaff returns internal members (admin + staff) for the settings page.
func (s *UserService) ListStaff(ctx context.Context, actor domain.Principal, search string) ([]domain.User, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can list users")
	}
	return s.users.ListStaff(ctx, search, 200)
}

// InviteStaff creates an internal account with an invitation status and a
// temporary password (the invitee uses password reset to set their own).
func (s *UserService) InviteStaff(ctx context.Context, actor domain.Principal, input domain.InviteInput) (*domain.User, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can invite users")
	}
	input.Sanitize()
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := validateName(input.Name); err != nil {
		return nil, err
	}
	if input.Role == "" {
		input.Role = domain.RoleStaff
	}
	if !input.Role.Valid() || !input.Role.StaffRole() {
		return nil, domain.Validationf("role must be one of: admin, staff")
	}
	if input.Department == "" {
		input.Department = "General Operations"
	}
	if len(input.Department) > 100 {
		return nil, domain.Validationf("department must not exceed 100 characters")
	}

	tmp, err := temporaryPassword()
	if err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(tmp)
	if err != nil {
		return nil, err
	}
	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         input.Role,
		Status:       domain.UserInvited,
		Department:   input.Department,
	}
	if err := s.users.Create(ctx, &user); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return nil, domain.Conflictf("a user with this email already exists")
		}
		return nil, err
	}
	return &user, nil
}

// GetUser returns a single user (staff only).
func (s *UserService) GetUser(ctx context.Context, actor domain.Principal, id string) (*domain.User, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view users")
	}
	return s.users.GetByID(ctx, id)
}

// UpdateUser patches a user's profile as admin (role/status/department/name).
func (s *UserService) UpdateUser(ctx context.Context, actor domain.Principal, id string, patch domain.UpdateUserPatch) (*domain.User, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can update users")
	}
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(patch.Name) != "" {
		user.Name = strings.TrimSpace(patch.Name)
	}
	if patch.Role != "" {
		if !patch.Role.Valid() || !patch.Role.StaffRole() {
			return nil, domain.Validationf("role must be one of: admin, staff")
		}
		user.Role = patch.Role
	}
	if patch.Status != "" {
		if !patch.Status.Valid() {
			return nil, domain.Validationf("invalid status %q", patch.Status)
		}
		user.Status = patch.Status
	}
	if strings.TrimSpace(patch.Department) != "" {
		user.Department = strings.TrimSpace(patch.Department)
	}
	return s.users.UpdateAdmin(ctx, user)
}

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

func temporaryPassword() (string, error) {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return "Invite-" + base64.RawURLEncoding.EncodeToString(buffer), nil
}

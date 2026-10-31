package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

const testSecret = "test-secret-that-is-longer-than-32-characters!!"

func newTestService(store *fakeStore) *AuthService {
	tokens := auth.NewTokenManager(testSecret, "brightloop-auth", time.Hour)
	return NewAuthService(store, store, tokens, 30*time.Minute)
}

func validRegister() domain.RegisterInput {
	return domain.RegisterInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "StrongPass1!",
		Role:     domain.RoleCustomer,
	}
}

func invalidPasswords() []string {
	return []string{"short", "alllowercase1!", "ALLUPPERCASE1!", "NoSpecial123", "NoUpperlower1", ""}
}

func TestRegisterCreatesCustomer(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)

	user, token, err := svc.Register(context.Background(), validRegister())
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	clean := user
	if clean.PasswordHash != "" && auth.VerifyPassword(clean.PasswordHash, "StrongPass1!") == false {
		t.Fatal("password hash does not match")
	}
	if user.Email != "jane@example.com" {
		t.Fatalf("email = %q", user.Email)
	}
	if user.Role != domain.RoleCustomer {
		t.Fatalf("role = %q, want customer", user.Role)
	}
	if token == "" {
		t.Fatal("expected a token")
	}
	if len(store.created) != 1 {
		t.Fatalf("created users = %d, want 1", len(store.created))
	}
}

func TestRegisterRejectsWeakPasswords(t *testing.T) {
	for _, password := range invalidPasswords() {
		store := newFakeStore()
		svc := newTestService(store)
		input := validRegister()
		input.Password = password
		if _, _, err := svc.Register(context.Background(), input); !errors.Is(err, domain.ErrValidation) {
			t.Errorf("password %q: err = %v, want validation", password, err)
		}
	}
}

func TestRegisterRejectsNonCustomerRole(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	input := validRegister()
	input.Role = domain.RoleAdmin
	if _, _, err := svc.Register(context.Background(), input); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("err = %v, want forbidden", err)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if _, _, err := svc.Register(context.Background(), validRegister()); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, _, err := svc.Register(context.Background(), validRegister()); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("second register err = %v, want conflict", err)
	}
}

func TestLoginSuccessAndWrongPassword(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if _, _, err := svc.Register(context.Background(), validRegister()); err != nil {
		t.Fatalf("register: %v", err)
	}

	user, token, err := svc.Login(context.Background(), domain.LoginInput{Email: "jane@example.com", Password: "StrongPass1!"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if token == "" || user.ID == "" {
		t.Fatal("expected token and user")
	}

	if _, _, err := svc.Login(context.Background(), domain.LoginInput{Email: "jane@example.com", Password: "WrongPass1!"}); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("wrong password err = %v, want unauthorized", err)
	}
	if _, _, err := svc.Login(context.Background(), domain.LoginInput{Email: "ghost@example.com", Password: "StrongPass1!"}); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("unknown email err = %v, want unauthorized", err)
	}
}

func TestValidateAfterLogoutRejectsToken(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if _, _, err := svc.Register(context.Background(), validRegister()); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, token, err := svc.Login(context.Background(), domain.LoginInput{Email: "jane@example.com", Password: "StrongPass1!"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if _, err := svc.Validate(context.Background(), token); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if err := svc.Logout(context.Background(), token); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := svc.Validate(context.Background(), token); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("validate after logout err = %v, want unauthorized", err)
	}
}

func TestValidateRejectsGarbageToken(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if _, err := svc.Validate(context.Background(), "not-a-jwt"); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("err = %v, want unauthorized", err)
	}
}

func TestResetPasswordCreatesTokenForExistingUser(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if _, _, err := svc.Register(context.Background(), validRegister()); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.ResetPassword(context.Background(), "jane@example.com"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if len(store.resets) != 1 {
		t.Fatalf("reset tokens = %d, want 1", len(store.resets))
	}
}

func TestResetPasswordIsGenericForUnknownEmail(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if err := svc.ResetPassword(context.Background(), "nobody@example.com"); err != nil {
		t.Fatalf("reset unknown email err = %v, want nil", err)
	}
	if len(store.resets) != 0 {
		t.Fatalf("reset tokens = %d, want 0 (no enumeration)", len(store.resets))
	}
}

func TestResetPasswordRejectsInvalidEmail(t *testing.T) {
	store := newFakeStore()
	svc := newTestService(store)
	if err := svc.ResetPassword(context.Background(), "not-an-email"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("err = %v, want validation", err)
	}
}

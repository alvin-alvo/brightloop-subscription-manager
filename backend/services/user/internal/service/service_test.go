package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

func newService(store *fakeStore) *UserService {
	return NewUserService(store)
}

func adminActor() domain.Principal {
	return domain.Principal{UserID: "usr_000", Role: domain.RoleAdmin}
}

func customerActor() domain.Principal {
	return domain.Principal{UserID: "usr_001", Role: domain.RoleCustomer}
}

func TestMeRequiresIdentity(t *testing.T) {
	svc := newService(newFakeStore())
	if _, err := svc.Me(context.Background(), domain.Principal{}); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("err = %v, want unauthorized", err)
	}
}

func TestListStaffForbidsNonStaff(t *testing.T) {
	svc := newService(newFakeStore())
	if _, err := svc.ListStaff(context.Background(), customerActor(), ""); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("err = %v, want forbidden", err)
	}
}

func TestInviteStaffCreatesInvitedUser(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	user, err := svc.InviteStaff(context.Background(), adminActor(), domain.InviteInput{
		Name: "Taylor Evans", Email: "taylor@brightloop.io", Role: domain.RoleStaff, Department: "Ops",
	})
	if err != nil {
		t.Fatalf("invite: %v", err)
	}
	if user.Role != domain.RoleStaff {
		t.Fatalf("role = %q", user.Role)
	}
	if user.Status != domain.UserInvited {
		t.Fatalf("status = %q, want invited", user.Status)
	}
	if user.PasswordHash == "" {
		t.Fatal("expected a password hash")
	}
}

func TestInviteRejectsCustomerRole(t *testing.T) {
	svc := newService(newFakeStore())
	if _, err := svc.InviteStaff(context.Background(), adminActor(), domain.InviteInput{
		Name: "A B", Email: "a@b.io", Role: domain.RoleCustomer, Department: "Ops",
	}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("err = %v, want validation", err)
	}
}

func TestInviteDuplicateEmail(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	input := domain.InviteInput{Name: "Taylor Evans", Email: "taylor@brightloop.io", Role: domain.RoleStaff}
	if _, err := svc.InviteStaff(context.Background(), adminActor(), input); err != nil {
		t.Fatalf("invite: %v", err)
	}
	if _, err := svc.InviteStaff(context.Background(), adminActor(), input); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("second invite err = %v, want conflict", err)
	}
}

func TestUpdateUserRole(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	created, err := svc.InviteStaff(context.Background(), adminActor(), domain.InviteInput{
		Name: "Taylor Evans", Email: "taylor@brightloop.io", Role: domain.RoleStaff,
	})
	if err != nil {
		t.Fatalf("invite: %v", err)
	}
	updated, err := svc.UpdateUser(context.Background(), adminActor(), created.ID, domain.UpdateUserPatch{Role: domain.RoleAdmin})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Role != domain.RoleAdmin {
		t.Fatalf("role = %q, want admin", updated.Role)
	}
}

func TestUpdateProfileKeepsEmail(t *testing.T) {
	store := newFakeStore()
	store.seed(domain.User{Name: "Alex", Email: "alex@brightloop.io", Role: domain.RoleAdmin, Status: domain.UserActive})
	admin := adminActor()
	admin.UserID = ""
	for id, u := range store.users {
		admin.UserID = id
		_ = u
	}
	svc := newService(store)
	updated, err := svc.UpdateProfile(context.Background(), admin, domain.UpdateProfilePatch{Name: "Alex Vance", Phone: "555"})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if updated.Email != "alex@brightloop.io" {
		t.Fatalf("email changed to %q", updated.Email)
	}
	if updated.Name != "Alex Vance" {
		t.Fatalf("name = %q", updated.Name)
	}
}

func TestDeleteAccountSetsInactive(t *testing.T) {
	store := newFakeStore()
	store.seed(domain.User{Name: "Alex", Email: "alex@brightloop.io", Role: domain.RoleAdmin, Status: domain.UserActive})
	admin := adminActor()
	for id := range store.users {
		admin.UserID = id
	}
	svc := newService(store)
	if err := svc.DeleteAccount(context.Background(), admin); err != nil {
		t.Fatalf("delete: %v", err)
	}
	got, _ := store.GetByID(context.Background(), admin.UserID)
	if got.Status != domain.UserInactive {
		t.Fatalf("status = %q, want inactive", got.Status)
	}
}

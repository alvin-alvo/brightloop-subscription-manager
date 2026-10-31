package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
)

func newService(store *fakeStore) *NotificationService {
	return NewNotificationService(store)
}

func TestListRequiresIdentity(t *testing.T) {
	svc := newService(newFakeStore())
	if _, _, err := svc.List(context.Background(), domain.Principal{}, 10, false); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("err = %v, want unauthorized", err)
	}
}

func TestSendValidates(t *testing.T) {
	svc := newService(newFakeStore())
	if _, err := svc.Send(context.Background(), domain.CreateNotificationInput{
		UserID: "", Channel: domain.ChannelInApp, Title: "Hi", Body: "Body",
	}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("missing user err = %v, want validation", err)
	}
	if _, err := svc.Send(context.Background(), domain.CreateNotificationInput{
		UserID: "u1", Channel: domain.ChannelInApp, Title: "", Body: "Body",
	}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("missing title err = %v, want validation", err)
	}
}

func TestSendCreatesUnread(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	n, err := svc.Send(context.Background(), domain.CreateNotificationInput{
		UserID: "u1", Channel: domain.ChannelInApp, Title: "Hello", Body: "World",
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if n.Read() {
		t.Fatal("expected unread")
	}
	items, count, err := svc.List(context.Background(), domain.Principal{UserID: "u1"}, 10, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || count != 1 {
		t.Fatalf("items = %d count = %d, want 1/1", len(items), count)
	}
}

func TestMarkAllReadClearsUnread(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	_, _ = svc.Send(context.Background(), domain.CreateNotificationInput{UserID: "u1", Channel: domain.ChannelInApp, Title: "A", Body: "B"})
	_, _ = svc.Send(context.Background(), domain.CreateNotificationInput{UserID: "u1", Channel: domain.ChannelInApp, Title: "C", Body: "D"})

	if err := svc.MarkAllRead(context.Background(), domain.Principal{UserID: "u1"}); err != nil {
		t.Fatalf("mark all: %v", err)
	}
	_, count, _ := svc.List(context.Background(), domain.Principal{UserID: "u1"}, 10, false)
	if count != 0 {
		t.Fatalf("unread = %d, want 0", count)
	}
}

func TestMarkReadUnknownIsNotFound(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	_, _ = svc.Send(context.Background(), domain.CreateNotificationInput{UserID: "u1", Channel: domain.ChannelInApp, Title: "A", Body: "B"})
	if err := svc.MarkRead(context.Background(), domain.Principal{UserID: "u1"}, "n_missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want not found", err)
	}
}

func TestCreateTemplateDuplicate(t *testing.T) {
	store := newFakeStore()
	svc := newService(store)
	inputs := []domain.CreateTemplateInput{
		{Code: "welcome", Channel: domain.ChannelInApp, Body: "Hi"},
		{Code: "WELCOME", Channel: domain.ChannelInApp, Body: "Ho"},
	}
	if _, err := svc.CreateTemplate(context.Background(), inputs[0]); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.CreateTemplate(context.Background(), inputs[1]); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("duplicate err = %v, want conflict", err)
	}
}

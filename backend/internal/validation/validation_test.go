package validation

import (
	"errors"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "owner@example.com"},
		{name: "valid with plus tag", value: "owner+tag@example.com"},
		{name: "missing at sign", value: "owner.example.com", wantErr: true},
		{name: "empty", value: "", wantErr: true},
		{name: "display name wrapper", value: "Owner <owner@example.com>", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Email(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Email(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, domain.ErrValidation) {
				t.Fatalf("Email(%q) error = %v, want domain.ErrValidation", tt.value, err)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "Str0ng!Pass"},
		{name: "exactly eight characters", value: "Ab1!Ab1!", wantErr: true},
		{name: "no uppercase", value: "str0ng!pass", wantErr: true},
		{name: "no lowercase", value: "STR0NG!PASS", wantErr: true},
		{name: "no special character", value: "Str0ngPass", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Password(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Password(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, domain.ErrValidation) {
				t.Fatalf("Password(%q) error = %v, want domain.ErrValidation", tt.value, err)
			}
		})
	}
}

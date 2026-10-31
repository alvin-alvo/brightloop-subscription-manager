package domain

import (
	"strings"
	"time"
)

// Channel mirrors the notification_channel enum in PostgreSQL.
type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelPush  Channel = "push"
	ChannelInApp Channel = "in_app"
)

// Valid reports whether c is a value accepted by the database enum.
func (c Channel) Valid() bool {
	switch c {
	case ChannelEmail, ChannelPush, ChannelInApp:
		return true
	default:
		return false
	}
}

// Notification is an in-app (or delivered) notification for a user.
type Notification struct {
	ID        string
	UserID    string
	Channel   Channel
	Title     string
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

// Read reports whether the notification has been read.
func (n Notification) Read() bool {
	return n.ReadAt != nil
}

// CreateNotificationInput is a validated notification.
type CreateNotificationInput struct {
	UserID  string
	Channel Channel
	Title   string
	Body    string
}

// Sanitize normalizes a notification input.
func (i *CreateNotificationInput) Sanitize() {
	i.UserID = strings.TrimSpace(i.UserID)
	i.Title = strings.TrimSpace(i.Title)
	i.Body = strings.TrimSpace(i.Body)
}

// Principal carries the identity resolved from the bearer token.
type Principal struct {
	UserID string
	Role   string
}

// NotificationTemplate is a reusable notification template.
type NotificationTemplate struct {
	ID        string
	Code      string
	Channel   Channel
	Subject   string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateTemplateInput is a validated template.
type CreateTemplateInput struct {
	Code    string
	Channel Channel
	Subject string
	Body    string
}

// Sanitize normalizes a template input.
func (i *CreateTemplateInput) Sanitize() {
	i.Code = strings.ToLower(strings.TrimSpace(i.Code))
	i.Subject = strings.TrimSpace(i.Subject)
	i.Body = strings.TrimSpace(i.Body)
}

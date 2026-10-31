// Package dto converts notification domain values to/from API shapes.
package dto

import (
	"strconv"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
)

// NotificationResponse is one notification for the bell dropdown.
type NotificationResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Read        bool   `json:"read"`
	CreatedAt   string `json:"createdAt"`
	DisplayTime string `json:"displayTime"`
}

// Notification converts a domain notification.
func Notification(n domain.Notification, now time.Time) NotificationResponse {
	return NotificationResponse{
		ID:          n.ID,
		Type:        string(n.Channel),
		Title:       n.Title,
		Body:        n.Body,
		Read:        n.Read(),
		CreatedAt:   timestamp(n.CreatedAt),
		DisplayTime: relativeTime(n.CreatedAt, now),
	}
}

// NotificationListResponse wraps the list for the bell dropdown.
type NotificationListResponse struct {
	Items       []NotificationResponse `json:"items"`
	UnreadCount int                    `json:"unreadCount"`
}

// UnreadCountResponse is the unread badge count.
type UnreadCountResponse struct {
	Count int `json:"count"`
}

// MessageResponse is a generic success payload.
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OKMessage builds a success message response.
func OKMessage(message string) MessageResponse {
	return MessageResponse{Success: true, Message: message}
}

// SendRequest is the notification creation payload.
type SendRequest struct {
	UserID  string `json:"userId"`
	Channel string `json:"channel"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

// ToInput converts to the domain input.
func (r SendRequest) ToInput() domain.CreateNotificationInput {
	return domain.CreateNotificationInput{
		UserID:  r.UserID,
		Channel: domain.Channel(r.Channel),
		Title:   r.Title,
		Body:    r.Body,
	}
}

// TemplateResponse is a stored notification template.
type TemplateResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Channel   string `json:"channel"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// Template converts a domain template.
func Template(t domain.NotificationTemplate) TemplateResponse {
	return TemplateResponse{
		ID:        t.ID,
		Code:      t.Code,
		Channel:   string(t.Channel),
		Subject:   t.Subject,
		Body:      t.Body,
		CreatedAt: timestamp(t.CreatedAt),
		UpdatedAt: timestamp(t.UpdatedAt),
	}
}

// CreateTemplateRequest is the template creation payload.
type CreateTemplateRequest struct {
	Code    string `json:"code"`
	Channel string `json:"channel"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// ToInput converts to the domain input.
func (r CreateTemplateRequest) ToInput() domain.CreateTemplateInput {
	return domain.CreateTemplateInput{
		Code:    r.Code,
		Channel: domain.Channel(r.Channel),
		Subject: r.Subject,
		Body:    r.Body,
	}
}

func timestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func relativeTime(value, now time.Time) string {
	value = value.UTC()
	now = now.UTC()
	diff := now.Sub(value)
	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		return strconv.Itoa(int(diff/time.Minute)) + "m ago"
	case diff < 24*time.Hour:
		return strconv.Itoa(int(diff/time.Hour)) + "h ago"
	case diff < 48*time.Hour:
		return "Yesterday"
	case diff < 7*24*time.Hour:
		return strconv.Itoa(int(diff/(24*time.Hour))) + "d ago"
	default:
		return value.Format("Jan 2, 2006")
	}
}

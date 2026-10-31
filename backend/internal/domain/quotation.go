package domain

import (
	"time"

	"github.com/google/uuid"
)

// QuotationTemplate is a reusable starting point for a quotation. It is valid
// for ValidDays after it is issued.
type QuotationTemplate struct {
	ID        uuid.UUID
	Name      string
	PlanID    uuid.UUID
	ValidDays int
	Lines     []QuotationTemplateLine
	CreatedAt time.Time
	UpdatedAt time.Time
}

// QuotationTemplateLine is a predefined product line on a template.
type QuotationTemplateLine struct {
	ID             uuid.UUID
	TemplateID     uuid.UUID
	ProductID      uuid.UUID
	Quantity       int
	UnitPriceCents int64
}

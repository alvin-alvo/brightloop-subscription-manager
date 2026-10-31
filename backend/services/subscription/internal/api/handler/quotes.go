package handler

import (
	"net/http"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/response"
	"github.com/gin-gonic/gin"
)

// Workspace handles GET /api/v1/subscriptions/workspace.
func (h *Handlers) Workspace(c *gin.Context) {
	workspace, err := h.Quotes.Workspace(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Workspace(workspace))
}

// ListQuotes handles GET /api/v1/subscriptions/quotes.
func (h *Handlers) ListQuotes(c *gin.Context) {
	quotes, err := h.Quotes.ListQuotes(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.QuoteResponse, 0, len(quotes))
	for _, q := range quotes {
		items = append(items, dto.Quote(q))
	}
	response.JSON(c, http.StatusOK, items)
}

// GetQuote handles GET /api/v1/subscriptions/quotes/:id.
func (h *Handlers) GetQuote(c *gin.Context) {
	quote, err := h.Quotes.GetQuote(c.Request.Context(), middleware.Actor(c), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Quote(*quote))
}

// CreateQuote handles POST /api/v1/subscriptions/quotes.
func (h *Handlers) CreateQuote(c *gin.Context) {
	var req dto.QuoteRequest
	if err := response.DecodeJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	q := req.ToDomain()
	quote, err := h.Quotes.CreateQuote(c.Request.Context(), middleware.Actor(c), &q)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Quote(*quote))
}

// UpdateQuoteStatus handles PATCH /api/v1/subscriptions/quotes/:id/status.
func (h *Handlers) UpdateQuoteStatus(c *gin.Context) {
	var body struct {
		Status string `json:"status"`
	}
	if err := response.DecodeJSON(c, &body); err != nil {
		response.Error(c, err)
		return
	}
	quote, err := h.Quotes.UpdateQuoteStatus(c.Request.Context(), middleware.Actor(c), c.Param("id"), body.Status)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Quote(*quote))
}

// DeleteQuote handles DELETE /api/v1/subscriptions/quotes/:id.
func (h *Handlers) DeleteQuote(c *gin.Context) {
	if err := h.Quotes.DeleteQuote(c.Request.Context(), middleware.Actor(c), c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("Quote deleted"))
}

// ListTemplates handles GET /api/v1/subscriptions/quotes/templates.
func (h *Handlers) ListTemplates(c *gin.Context) {
	templates, err := h.Quotes.ListTemplates(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.TemplateResponse, 0, len(templates))
	for _, t := range templates {
		items = append(items, dto.Template(t))
	}
	response.JSON(c, http.StatusOK, items)
}

// GetTemplate handles GET /api/v1/subscriptions/quotes/templates/:id.
func (h *Handlers) GetTemplate(c *gin.Context) {
	t, err := h.Quotes.GetTemplate(c.Request.Context(), middleware.Actor(c), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Template(*t))
}

// CreateTemplate handles POST /api/v1/subscriptions/quotes/templates.
func (h *Handlers) CreateTemplate(c *gin.Context) {
	var req dto.TemplateRequest
	if err := response.DecodeJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	t := req.ToDomain()
	created, err := h.Quotes.CreateTemplate(c.Request.Context(), middleware.Actor(c), &t)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Template(*created))
}

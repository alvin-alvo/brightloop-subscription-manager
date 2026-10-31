package handler

import (
	"net/http"
	"strconv"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/gin-gonic/gin"
)

// ListProducts handles GET /api/v1/subscriptions/products.
func (h *Handlers) ListProducts(c *gin.Context) {
	items, err := h.Catalog.ListProducts(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]dto.ProductResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, dto.Product(p))
	}
	response.JSON(c, http.StatusOK, resp)
}

// CreateProduct handles POST /api/v1/subscriptions/products.
func (h *Handlers) CreateProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := response.DecodeJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	p, err := h.Catalog.CreateProduct(c.Request.Context(), middleware.Actor(c), req.ToCreate())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Product(*p))
}

// UpdateProduct handles PUT /api/v1/subscriptions/products/:id.
func (h *Handlers) UpdateProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := response.DecodeJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	p, err := h.Catalog.UpdateProduct(c.Request.Context(), middleware.Actor(c), c.Param("id"), req.ToUpdate())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Product(*p))
}

// DeleteProduct handles DELETE /api/v1/subscriptions/products/:id.
func (h *Handlers) DeleteProduct(c *gin.Context) {
	if err := h.Catalog.DeleteProduct(c.Request.Context(), middleware.Actor(c), c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("Product deleted"))
}

// ListDiscounts handles GET /api/v1/subscriptions/discounts.
func (h *Handlers) ListDiscounts(c *gin.Context) {
	activeOnly := c.Query("activeOnly") == "true"
	items, err := h.Catalog.ListDiscounts(c.Request.Context(), middleware.Actor(c), activeOnly)
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]dto.DiscountResponse, 0, len(items))
	for _, d := range items {
		resp = append(resp, dto.Discount(d))
	}
	response.JSON(c, http.StatusOK, resp)
}

// CreateDiscount handles POST /api/v1/subscriptions/discounts.
func (h *Handlers) CreateDiscount(c *gin.Context) {
	var req dto.DiscountRequest
	if err := response.DecodeJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	d := req.ToDomain()
	created, err := h.Catalog.CreateDiscount(c.Request.Context(), middleware.Actor(c), &d)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Discount(*created))
}

// ToggleDiscount handles PATCH /api/v1/subscriptions/discounts/:id.
func (h *Handlers) ToggleDiscount(c *gin.Context) {
	active := c.Query("active") != "false"
	d, err := h.Catalog.ToggleDiscount(c.Request.Context(), middleware.Actor(c), c.Param("id"), active)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Discount(*d))
}

// OverdueInvoices handles GET /api/v1/subscriptions/invoices.
func (h *Handlers) OverdueInvoices(c *gin.Context) {
	limit, err := queryPositive(c, "limit", 100)
	if err != nil {
		response.Error(c, err)
		return
	}
	items, err := h.Catalog.OverdueInvoices(c.Request.Context(), middleware.Actor(c), true, limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]dto.OverdueInvoiceResponse, 0, len(items))
	for _, inv := range items {
		resp = append(resp, dto.OverdueInvoice(inv))
	}
	response.JSON(c, http.StatusOK, resp)
}

// CustomerInvoices handles GET /api/v1/subscriptions/invoices/customer.
func (h *Handlers) CustomerInvoices(c *gin.Context) {
	items, err := h.Catalog.CustomerInvoices(c.Request.Context(), middleware.Actor(c), c.Query("customerId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := make([]dto.CustomerInvoiceResponse, 0, len(items))
	for _, inv := range items {
		resp = append(resp, dto.CustomerInvoice(inv))
	}
	response.JSON(c, http.StatusOK, resp)
}

// CustomerOverview handles GET /api/v1/subscriptions/customers/:id/overview.
func (h *Handlers) CustomerOverview(c *gin.Context) {
	overview, err := h.Catalog.CustomerOverview(c.Request.Context(), middleware.Actor(c), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.CustomerOverview(overview))
}

// PayInvoice handles POST /api/v1/subscriptions/invoices/:id/pay.
func (h *Handlers) PayInvoice(c *gin.Context) {
	var body struct {
		Method string `json:"method"`
	}
	_ = response.DecodeJSON(c, &body)
	if body.Method == "" {
		body.Method = "card"
	}
	inv, err := h.Catalog.PayInvoice(c.Request.Context(), middleware.Actor(c), c.Param("id"), body.Method)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.CustomerInvoice(*inv))
}

func queryPositive(c *gin.Context, key string, fallback int) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return 0, domain.Validationf("%s must be a positive integer", key)
	}
	return v, nil
}

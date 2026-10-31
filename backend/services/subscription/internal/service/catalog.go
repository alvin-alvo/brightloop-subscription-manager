package service

import (
	"context"
	"errors"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/repository"
)

// CatalogService implements admin catalogue and AR/customer portal rules.
type CatalogService struct {
	catalog repository.CatalogRepository
	refs    repository.ReferenceRepository
}

// NewCatalogService wires the catalogue business rules.
func NewCatalogService(catalog repository.CatalogRepository, refs repository.ReferenceRepository) *CatalogService {
	return &CatalogService{catalog: catalog, refs: refs}
}

// --- Products --------------------------------------------------------------

func (s *CatalogService) ListProducts(ctx context.Context, actor domain.Actor) ([]domain.Product, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage the catalogue")
	}
	return s.catalog.ListProducts(ctx)
}

func (s *CatalogService) CreateProduct(ctx context.Context, actor domain.Actor, input domain.CreateProductInput) (*domain.Product, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage the catalogue")
	}
	if err := validateName("name", input.Name, true); err != nil {
		return nil, err
	}
	p := domain.Product{Name: input.Name, Description: input.Description, SKU: input.SKU, Active: true, Features: input.Features}
	if err := s.catalog.CreateProduct(ctx, &p); err != nil {
		return nil, wrapConflict(err, "a product with this name or SKU already exists")
	}
	return &p, nil
}

func (s *CatalogService) UpdateProduct(ctx context.Context, actor domain.Actor, id string, input domain.UpdateProductInput) (*domain.Product, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage the catalogue")
	}
	active := true
	if input.Active != nil {
		active = *input.Active
	}
	p := domain.Product{ID: id, Name: input.Name, Description: input.Description, SKU: input.SKU, Active: active, Features: input.Features}
	if err := s.catalog.UpdateProduct(ctx, &p); err != nil {
		return nil, err
	}
	return s.catalog.GetProduct(ctx, id)
}

func (s *CatalogService) DeleteProduct(ctx context.Context, actor domain.Actor, id string) error {
	if !actor.IsStaff() {
		return domain.Forbiddenf("only staff can manage the catalogue")
	}
	return s.catalog.DeleteProduct(ctx, id)
}

// --- Discounts -------------------------------------------------------------

func (s *CatalogService) ListDiscounts(ctx context.Context, actor domain.Actor, activeOnly bool) ([]domain.Discount, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage discounts")
	}
	return s.catalog.ListDiscounts(ctx, activeOnly)
}

func (s *CatalogService) CreateDiscount(ctx context.Context, actor domain.Actor, d *domain.Discount) (*domain.Discount, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage discounts")
	}
	if err := validateDiscount(d); err != nil {
		return nil, err
	}
	if err := s.catalog.CreateDiscount(ctx, d); err != nil {
		return nil, wrapConflict(err, "a discount with this code already exists")
	}
	return d, nil
}

func (s *CatalogService) ToggleDiscount(ctx context.Context, actor domain.Actor, id string, active bool) (*domain.Discount, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage discounts")
	}
	if err := s.catalog.UpdateDiscount(ctx, &domain.Discount{ID: id, Active: active}); err != nil {
		return nil, err
	}
	return s.catalog.GetDiscount(ctx, id)
}

// --- AR & customer portal --------------------------------------------------

func (s *CatalogService) OverdueInvoices(ctx context.Context, actor domain.Actor, overdueOnly bool, limit int) ([]domain.OverdueInvoice, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view receivables")
	}
	return s.catalog.ListOverdueInvoices(ctx, overdueOnly, limit)
}

func (s *CatalogService) CustomerInvoices(ctx context.Context, actor domain.Actor, customerID string) ([]domain.CustomerInvoice, error) {
	id, err := s.resolveCustomer(ctx, actor, customerID)
	if err != nil {
		return nil, err
	}
	return s.catalog.CustomerInvoices(ctx, id)
}

func (s *CatalogService) CustomerOverview(ctx context.Context, actor domain.Actor, customerID string) (domain.CustomerOverview, error) {
	id, err := s.resolveCustomer(ctx, actor, customerID)
	if err != nil {
		return domain.CustomerOverview{}, err
	}
	return s.catalog.CustomerOverview(ctx, id)
}

// PayInvoice settles a customer's open invoice. Customers may only settle their
// own invoices; staff may settle any.
func (s *CatalogService) PayInvoice(ctx context.Context, actor domain.Actor, invoiceID, method string) (*domain.CustomerInvoice, error) {
	inv, err := s.catalog.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	customerID := inv.CustomerID
	if !actor.IsStaff() {
		resolved, err := s.resolveCustomer(ctx, actor, "")
		if err != nil {
			return nil, err
		}
		if resolved != customerID {
			return nil, domain.Forbiddenf("cannot pay another customer's invoice")
		}
	}
	return s.catalog.RecordPayment(ctx, invoiceID, customerID, method)
}

// resolveCustomer returns the target customer for the actor. Staff may pass a
// specific customerId; customers are always scoped to their own account.
func (s *CatalogService) resolveCustomer(ctx context.Context, actor domain.Actor, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "me" {
		requested = ""
	}
	if !actor.IsStaff() && requested == "" {
		if actor.CustomerID != "" {
			return actor.CustomerID, nil
		}
		if actor.UserID == "" {
			return "", domain.Forbiddenf("customer identity is not available")
		}
		c, err := s.refs.GetCustomerByUserID(ctx, actor.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return "", domain.Forbiddenf("no customer profile is linked to this account")
			}
			return "", err
		}
		return c.ID, nil
	}
	if actor.IsStaff() {
		if requested != "" {
			return requested, nil
		}
		if actor.CustomerID != "" {
			return actor.CustomerID, nil
		}
		if actor.UserID != "" {
			if c, err := s.refs.GetCustomerByUserID(ctx, actor.UserID); err == nil {
				return c.ID, nil
			}
		}
		return "", domain.Validationf("customerId is required")
	}
	return "", domain.Forbiddenf("cannot access another customer's data")
}

func wrapConflict(err error, message string) error {
	if errors.Is(err, domain.ErrConflict) {
		return domain.Conflictf("a record with these details already exists")
	}
	_ = message
	return err
}

func validateDiscount(d *domain.Discount) error {
	if d == nil || d.Code == "" {
		return domain.Validationf("discount code is required")
	}
	if d.Kind != domain.DiscountPercentage && d.Kind != domain.DiscountFixedAmount {
		return domain.Validationf("discount kind must be percentage or fixed_amount")
	}
	if d.Duration != domain.DiscountOnce && d.Duration != domain.DiscountRepeating && d.Duration != domain.DiscountForever {
		return domain.Validationf("discount duration must be once, repeating or forever")
	}
	if d.Kind == domain.DiscountPercentage && (d.PercentBasisPoints < 0 || d.PercentBasisPoints > 10000) {
		return domain.Validationf("percentage discount must be between 0 and 10000 basis points")
	}
	return nil
}

func validateName(field, value string, required bool) error {
	value = strings.TrimSpace(value)
	if required && value == "" {
		return domain.Validationf("%s is required", field)
	}
	if len(value) > 200 {
		return domain.Validationf("%s must not exceed 200 characters", field)
	}
	return nil
}

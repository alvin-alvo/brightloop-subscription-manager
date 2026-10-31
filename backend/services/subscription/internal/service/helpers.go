package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func utcDate(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// addMonthsClamped adds calendar months while clamping to the last valid day
// (Jan 31 + 1 month => Feb 28/29).
func addMonthsClamped(date time.Time, months int) time.Time {
	if months == 0 {
		return date
	}
	year, month, day := date.Date()
	firstOfTarget := time.Date(year, month+time.Month(months), 1, 0, 0, 0, 0, date.Location())
	lastDay := time.Date(firstOfTarget.Year(), firstOfTarget.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(firstOfTarget.Year(), firstOfTarget.Month(), day, 0, 0, 0, 0, date.Location())
}

func roundPercent(numerator, denominator int64) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator) * 100
}

func percentChange(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(previous) * 100
}

func roundCents(value float64) int64 {
	return int64(math.Round(value))
}

func upperCurrency(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) != 3 {
		return ""
	}
	return value
}

func invoiceNumber(now time.Time) string {
	buffer := make([]byte, 4)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("INV-%s", now.UTC().Format("20060102150405"))
	}
	return fmt.Sprintf("INV-%s-%s", now.UTC().Format("20060102"), strings.ToUpper(hex.EncodeToString(buffer)))
}

func discountAmountCents(discount domain.Discount, subtotal int64) int64 {
	var amount int64
	switch discount.Kind {
	case domain.DiscountPercentage:
		bp := int64(discount.PercentBasisPoints)
		if bp < 0 {
			bp = 0
		}
		if bp > 10000 {
			bp = 10000
		}
		amount = (subtotal*bp + 5000) / 10000
	case domain.DiscountFixedAmount:
		amount = discount.AmountCents
	}
	if amount < 0 {
		return 0
	}
	if amount > subtotal {
		return subtotal
	}
	return amount
}

func discountUsable(discount domain.Discount, currency string, now time.Time) error {
	if !discount.Active {
		return domain.Validationf("discount %q is not active", discount.Code)
	}
	if discount.MaxRedemptions != nil && discount.TimesRedeemed >= *discount.MaxRedemptions {
		return domain.Validationf("discount %q has reached its redemption limit", discount.Code)
	}
	if discount.StartsAt != nil && now.Before(*discount.StartsAt) {
		return domain.Validationf("discount %q is not active yet", discount.Code)
	}
	if discount.EndsAt != nil && now.After(*discount.EndsAt) {
		return domain.Validationf("discount %q has expired", discount.Code)
	}
	if discount.Currency != "" && currency != "" && !strings.EqualFold(discount.Currency, currency) {
		return domain.Validationf("discount %q does not apply to currency %s", discount.Code, currency)
	}
	return nil
}

func validatePeriod(start, end time.Time) error {
	if end.Before(start) || end.Equal(start) {
		return domain.Validationf("currentPeriodEnd must be after currentPeriodStart")
	}
	return nil
}

func pageCount(total, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func statusCounts(a domain.Analytics) []domain.StatusBreakdown {
	counts := []struct {
		status string
		count  int
	}{
		{string(domain.StatusActive), a.ActiveSubscriptions},
		{string(domain.StatusTrialing), a.TrialingSubscriptions},
		{string(domain.StatusPastDue), a.PastDueSubscriptions},
		{string(domain.StatusPaused), a.PausedSubscriptions},
		{string(domain.StatusUnpaid), a.UnpaidSubscriptions},
		{string(domain.StatusCanceled), a.CanceledSubscriptions},
	}
	total := 0
	for _, item := range counts {
		total += item.count
	}
	result := make([]domain.StatusBreakdown, 0, len(counts))
	for _, item := range counts {
		percent := 0.0
		if total > 0 {
			percent = float64(item.count) / float64(total) * 100
		}
		result = append(result, domain.StatusBreakdown{
			Status:  item.status,
			Count:   item.count,
			Percent: percent,
		})
	}
	return result
}

func computeAnalytics(a domain.Analytics) domain.ComputedAnalytics {
	computed := domain.ComputedAnalytics{Analytics: a}
	computed.ARR = a.MRRCents * 12
	active := int64(a.ActiveSubscriptions)
	if active > 0 {
		computed.ARPUCents = (a.MRRCents + active/2) / active
	}
	computed.MRRChangePercent = percentChange(a.MRRCents, a.PreviousMRRCents)
	computed.NetNewSubscriptions = a.NewSubscriptions - a.CanceledThisMonth

	previousActive := a.ActiveSubscriptions - a.NewSubscriptions + a.CanceledThisMonth
	if previousActive < 0 {
		previousActive = 0
	}
	computed.ActiveChangePercent = percentChange(int64(a.ActiveSubscriptions), int64(previousActive))

	previousCustomers := a.TotalCustomers - a.NewCustomers
	if previousCustomers < 0 {
		previousCustomers = 0
	}
	computed.CustomersChangePercent = percentChange(int64(a.TotalCustomers), int64(previousCustomers))

	if a.PreviousMRRCents > 0 {
		computed.NetRetentionRate = float64(a.NetRetainedMRRCents) / float64(a.PreviousMRRCents) * 100
	} else if a.MRRCents > 0 {
		// No prior-period base: report the current book as fully retained.
		computed.NetRetentionRate = 100
	}
	if previousActive > 0 {
		computed.ChurnRate = float64(a.CanceledThisMonth) / float64(previousActive) * 100
	}
	return computed
}

package controllers

import (
	"sync"
	"time"

	"github.com/ericklima-ca/rinhador/models"
	"github.com/ericklima-ca/rinhador/services"
	"github.com/gofiber/fiber/v3"
)

// Service used: Default or Fallback
type Service string

const (
	Default  Service = "default"
	Fallback Service = "fallback"
)

// Persistent in-memory DB instance
var db = (&InMemoryDB{}).New()

type Transaction struct {
	CorrelationID string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
	Service       Service   `json:"service"`
}

type InMemoryDB struct {
	Transactions []Transaction
	mu           sync.Mutex
}

func (db *InMemoryDB) New() *InMemoryDB {
	return &InMemoryDB{
		Transactions: []Transaction{},
		mu:           sync.Mutex{},
	}
}

func (db *InMemoryDB) GetSummary(from, to time.Time) models.Summary {
	db.mu.Lock()
	defer db.mu.Unlock()

	var summary models.Summary
	for _, tx := range db.Transactions {
		// Filter by date range if provided
		if !from.IsZero() && tx.RequestedAt.Before(from) {
			continue
		}
		if !to.IsZero() && tx.RequestedAt.After(to) {
			continue
		}
		switch tx.Service {
		case Default:
			summary.Default.TotalRequests++
			summary.Default.TotalAmount += tx.Amount
		case Fallback:
			summary.Fallback.TotalRequests++
			summary.Fallback.TotalAmount += tx.Amount
		}
	}
	return summary
}

func (db *InMemoryDB) AddTransaction(correlationID string, amount float64, now time.Time, service Service) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Transactions = append(db.Transactions, Transaction{
		CorrelationID: correlationID,
		Amount:        amount,
		Service:       service,
		RequestedAt:   now,
	})
}

func Payments(c fiber.Ctx) error {
	var payment models.Payment
	now := time.Now().UTC()

	// Parse the JSON body into the Payment struct
	if err := c.Bind().JSON(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment payload",
		})
	}

	var service Service
	var err error

	// Try processing with the default service
	if err = services.ProcessPayment(payment.CorrelationID, payment.Amount); err != nil {
		// Fallback if default fails
		if err = services.ProcessPaymentFallback(payment.CorrelationID, payment.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process payment",
			})
		}
		service = Fallback
		db.AddTransaction(payment.CorrelationID, payment.Amount, now, service)
	} else {
		service = Default
		db.AddTransaction(payment.CorrelationID, payment.Amount, now, service)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Payment processed successfully",
		"correlationId": payment.CorrelationID,
		"service":       service,
	})
}

type PaymentsSummaryParams struct {
	From string `query:"from"`
	To   string `query:"to"`
}

func PaymentsSummary(c fiber.Ctx) error {
	var params PaymentsSummaryParams
	if err := c.Bind().Query(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	parseTime := func(value, field string) (time.Time, error) {
		if value == "" {
			return time.Time{}, nil
		}
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid '"+field+"' date format")
		}
		return t, nil
	}

	from, err := parseTime(params.From, "from")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	to, err := parseTime(params.To, "to")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	summary := db.GetSummary(from, to)
	return c.Status(fiber.StatusOK).JSON(summary)
}

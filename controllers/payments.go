package controllers

import (
	"fmt"

	"github.com/ericklima-ca/rinhador/models"
	"github.com/gofiber/fiber/v3"
)

// ## Payments

// > Principal endpoint que recebe requisições de pagamentos a serem processados.

// ```
// POST /payments
// {
//     "correlationId": "4a7901b8-7d26-4d9d-aa19-4dc1c7cf60b3",
//     "amount": 19.90
// }

// HTTP 2XX
// Qualquer coisa
// ```

// ### requisição

// - `correlationId` é um campo obrigatório e único do tipo UUID.
// - `amount` é um campo obrigatório do tipo decimal.

// ### resposta

// - Qualquer resposta na faixa 2XX (200, 201, 202, etc) é válida. O corpo da resposta não será validado – pode ser qualquer coisa ou até vazio.

func Payments(c fiber.Ctx) error {
	payment := new(models.Payment)

	// Parse the JSON body into the Payment struct
	if err := c.Bind().Body(payment); err != nil {
		return err
	}

	fmt.Printf("Processing payment with CorrelationID: %s and Amount: %.2f\n", payment.CorrelationID, payment.Amount)

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment processed successfully",
	})
}

// GET /payments-summary?from=2020-07-10T12:34:56.000Z&to=2020-07-10T12:35:56.000Z

// HTTP 200 - Ok
// {
//     "default" : {
//         "totalRequests": 43236,
//         "totalAmount": 415542345.98
//     },
//     "fallback" : {
//         "totalRequests": 423545,
//         "totalAmount": 329347.34
//     }
// }

type PaymentsSummaryParams struct {
	From string `query:"from"`
	To   string `query:"to"`
}

func PaymentsSummary(c fiber.Ctx) error {
	params := new(PaymentsSummaryParams)

	if err := c.Bind().Query(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}
	fmt.Printf("Fetching payment summary from %s to %s\n", params.From, params.To)

	summary := models.Summary{
		Default: struct {
			TotalRequests int     `json:"totalRequests" validate:"required"`
			TotalAmount   float64 `json:"totalAmount" validate:"required"`
		}{
			TotalRequests: 43236,
			TotalAmount:   415542345.98,
		},
		Fallback: struct {
			TotalRequests int     `json:"totalRequests" validate:"required"`
			TotalAmount   float64 `json:"totalAmount" validate:"required"`
		}{
			TotalRequests: 423545,
			TotalAmount:   329347.34,
		},
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

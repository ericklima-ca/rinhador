package models

type Payment struct {
	CorrelationID string  `json:"correlationId" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

type Summary struct {
	Default struct {
		TotalRequests int     `json:"totalRequests" validate:"required"`
		TotalAmount   float64 `json:"totalAmount" validate:"required"`
	} `json:"default" validate:"required"`
	Fallback struct {
		TotalRequests int     `json:"totalRequests" validate:"required"`
		TotalAmount   float64 `json:"totalAmount" validate:"required"`
	} `json:"fallback" validate:"required"`
}

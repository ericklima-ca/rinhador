package services

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var HTTPClient = &fasthttp.Client{
	MaxResponseBodySize: 4 * 1024 * 1024,
	ReadTimeout:         30 * time.Second,
	WriteTimeout:        30 * time.Second,

	Dial: (&fasthttp.TCPDialer{
		Concurrency: 300 * 2,
	}).Dial,

	MaxConnsPerHost:               300 * 2,
	MaxIdleConnDuration:           90 * time.Second,
	MaxConnDuration:               0,
	MaxConnWaitTimeout:            5 * time.Second,
	DisablePathNormalizing:        true,
	DisableHeaderNamesNormalizing: true,
}

func ProcessPayment(correlationID string, amount float64) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	requestBody := fmt.Sprintf(`{"correlationId": "%s", "amount": %.2f}`, correlationID, amount)
	req.SetRequestURI("http://localhost:8001/payments")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBodyString(requestBody)

	if err := HTTPClient.Do(req, resp); err != nil {
		// print error
		fmt.Printf("Error processing payment: %v\n", err)
		return fmt.Errorf("default payment services failed: %v", err)

	}
	if resp.StatusCode() != fasthttp.StatusOK {
		// print error
		fmt.Printf("Error processing payment, status code: %d\n", resp.StatusCode())
		return fmt.Errorf("default payment service returned status code: %d", resp.StatusCode())
	}
	return nil
}

func ProcessPaymentFallback(correlationID string, amount float64) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	requestBody := fmt.Sprintf(`{"correlationId": "%s", "amount": %.2f}`, correlationID, amount)
	req.SetRequestURI("http://localhost:8002/payments")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBodyString(requestBody)

	if err := HTTPClient.Do(req, resp); err != nil {
		// print error
		fmt.Printf("Error processing fallback payment: %v\n", err)
		return fmt.Errorf("fallback payment service failed: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		// print error
		fmt.Printf("Error processing fallback payment, status code: %d\n", resp.StatusCode())
		return fmt.Errorf("fallback payment service returned status code: %d", resp.StatusCode())
	}
	return nil
}

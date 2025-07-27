package main

import (
	"log"

	"github.com/ericklima-ca/rinhador/controllers"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	// "github.com/gofiber/fiber/v3/middleware/logger"
)

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	// Initialize a new Fiber app
	// Setup your validator in the config
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})

	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${latency} ${status} - ${method} ${path}/${queryParams}\n",
	// }))

	app.Post("/payments", controllers.Payments)
	app.Get("/payments-summary", controllers.PaymentsSummary)

	// Start the server on port 9999
	log.Fatal(app.Listen(":9999"))
}

package main

import (
    "log"
	"os"

    "github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/go-playground/validator/v10"
)

var (
	app	   *fiber.App
	PORT	string
	API_KEY	string
	validate = validator.New()
)

type EmailRequest struct {
	ApiKey string `json:"apiKey" form:"apiKey" validate:"required"`
	Username   string `json:"username" form:"username" validate:"required"`
	To	   string `json:"to" form:"to" validate:"required,email"`
	Subject string `json:"subject" form:"subject" validate:"required"`
	HTML   string `json:"html" form:"html" validate:"required"`
	Text   string `json:"text" form:"text" validate:"required"`
	Author string `json:"author" form:"author" validate:"required"`
}

func init() {

	// Load environment variables from a .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	PORT = os.Getenv("PORT")
	API_KEY = os.Getenv("API_KEY")
}

func main() {
	// Initialize the Fiber app with default settings
	app = fiber.New()

    // Define a route for the GET method on the root path '/'
    app.Get("/", func(c fiber.Ctx) error {
        // Send a string response to the client
        return c.JSON(fiber.Map{ 
			"message": "Welcome to MailRelay API!",
			"version": "1.0",
		})
    })

	app.Post("/send", func(c fiber.Ctx) error {
		emailReq := new(EmailRequest)

		if err := c.Bind().Body(emailReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if err := validate.Struct(emailReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Validation failed",
				"fields": err.Error(),
			})
		}

		if emailReq.ApiKey != API_KEY {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Email received",
			"data":    emailReq,
		})
	})

    // Start the server on port 3000
    log.Fatal(app.Listen(":" + PORT))
}
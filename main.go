package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
)

var (
	app			*fiber.App
	PORT		string
	API_KEY		string
	validate = validator.New()
	MAIL_HOST	string
	MAIL_DOMAIN	string
	MAIL_PORT	int
	MAIL_USER	string
	MAIL_PASS	string
	MAIL_SSL	bool
)

type EmailRequest struct {
	ApiKey		string `json:"apiKey" form:"apiKey" validate:"required"`
	Username	string `json:"username" form:"username" validate:"required"`
	To			string `json:"to" form:"to" validate:"required,email"`
	Subject		string `json:"subject" form:"subject" validate:"required"`
	HTML		string `json:"html" form:"html" validate:"required"`
	Text		string `json:"text" form:"text" validate:"required"`
	Author		string `json:"author" form:"author" validate:"required"`
}

func init() {

	// Load environment variables from a .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	PORT = os.Getenv("PORT")
	API_KEY = os.Getenv("API_KEY")
	MAIL_HOST = os.Getenv("MAIL_HOST")
	MAIL_DOMAIN = os.Getenv("MAIL_DOMAIN")
	mailPort := os.Getenv("MAIL_PORT")
	MAIL_PORT = 587
	if mailPort != "" {
		if port, err := strconv.Atoi(mailPort); err == nil {
			MAIL_PORT = port
		}
	}
	MAIL_USER = os.Getenv("MAIL_USER")
	MAIL_PASS = os.Getenv("MAIL_PASS")
	MAIL_SSL = os.Getenv("MAIL_SSL") == "true"

	if PORT == "" {
		PORT = "3000"
	}
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

		m := gomail.NewMessage()
		m.SetHeader("From", emailReq.Author + " <" + emailReq.Username + "@" + MAIL_DOMAIN + ">")
		m.SetHeader("To", emailReq.To)
		m.SetHeader("Subject", emailReq.Subject)
		m.SetBody("text/plain", emailReq.Text)
		m.AddAlternative("text/html", emailReq.HTML)

		d := gomail.NewDialer(MAIL_HOST, MAIL_PORT, MAIL_USER, MAIL_PASS)

		d.SSL = MAIL_SSL;

		if err := d.DialAndSend(m); err != nil {
			log.Println("Error sending email:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to send email",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email sent successfully",
		})
	})

    // Start the server on port 3000
    log.Fatal(app.Listen(":" + PORT))
}
package main

import (
    "log"
	"os"

    "github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

var app *fiber.App

var PORT string
var API_KEY string

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
		c.Accepts("application/json")
		c.Accepts("application/x-www-form-urlencoded")
		
		// Create a new email request struct
		emailReq := new(EmailRequest)
		
		// Parse body into the struct
		if err := c.Bind().Body(emailReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse request body",
			})
		}
		
		// Now you can access the fields
		apiKey := emailReq.APIKey
		title := emailReq.Title
		html := emailReq.HTML
		text := emailReq.Text
		author := emailReq.Author

		return c.JSON(fiber.Map{
			"apiKey": apiKey,
			"title":  title,
			"html":   html,
			"text":   text,
			"author": author,
		});
	})

    // Start the server on port 3000
    log.Fatal(app.Listen(":" + PORT))
}

type EmailRequest struct {
    APIKey string `json:"apiKey" form:"apiKey"`
    Title  string `json:"title" form:"title"`
    HTML   string `json:"html" form:"html"`
    Text   string `json:"text" form:"text"`
    Author string `json:"author" form:"author"`
}
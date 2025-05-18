package main

import (
    "log"
	"os"

    "github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

var PORT string
var app *fiber.App

func init() {

	// Load environment variables from a .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	PORT = os.Getenv("PORT")

	log.Println(PORT);
}

func main() {
	// Initialize the Fiber app with default settings
	app := fiber.New()

    // Define a route for the GET method on the root path '/'
    app.Get("/", func(c fiber.Ctx) error {
        // Send a string response to the client
        return c.JSON(fiber.Map{ 
			"message": "Welcome to MailRelay API!",
			"version": "1.0",
		})
    })

    // Start the server on port 3000
    log.Fatal(app.Listen(":" + PORT))
}
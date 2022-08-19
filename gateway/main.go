package main

import (
	routes "gateway/handlers"
	"gateway/services/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.InitDB()
	app := fiber.New()
	// api := app.Group("/api", logger.New())
	api := app.Group("/api")

	routes.AssignV1Handlers(api)

	app.Listen(":3000")
}

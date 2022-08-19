package handlers

import "github.com/gofiber/fiber/v2"

func AssignHelloHandlers(r fiber.Router) {
	r.Get("/hello", getHello)
}

func getHello(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

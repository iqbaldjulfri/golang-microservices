package routes

import (
	"gateway/handlers/v1"

	"github.com/gofiber/fiber/v2"
)

func AssignV1Handlers(api fiber.Router) {
	v1 := api.Group("v1")

	handlers.AssignHelloHandlers(v1)
	handlers.AssignUsersHandlers(v1)
}

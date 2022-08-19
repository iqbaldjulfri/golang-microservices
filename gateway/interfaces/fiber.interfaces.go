package interfaces

import "github.com/gofiber/fiber/v2"

type FiberBodyParser interface {
	BodyParser(interface{}) error
}

type FiberJSONSender interface {
	JSON(interface{}) error
}

type FiberLocalsGetterSetter interface {
	Locals(string, ...interface{}) interface{}
}

type FiberNextRunner interface {
	Next() error
}

type FiberStatusSetter interface {
	Status(int) *fiber.Ctx
}

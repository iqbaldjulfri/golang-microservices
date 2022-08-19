package middlewares

import (
	"gateway/interfaces"
	"gateway/utils"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type bodyBuilder func() interface{}

type requestBodyParser interface {
	interfaces.FiberBodyParser
	interfaces.FiberJSONSender
	interfaces.FiberLocalsGetterSetter
	interfaces.FiberNextRunner
	interfaces.FiberStatusSetter
}

func ValidateBodyFnFactory(f bodyBuilder) fiber.Handler {
	body := f()

	return func(c *fiber.Ctx) error {
		return validateBody(c, body)
	}
}

func validateBody(c requestBodyParser, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err, nil)
	}

	if err := validate.Struct(body); err != nil {
		m := fiber.Map{}
		for _, err := range err.(validator.ValidationErrors) {
			m[strings.ToLower(err.Field())] = err.Tag()
		}
		return utils.JSONStatus(c, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, m)
	}

	c.Locals("body", body)

	return c.Next()
}

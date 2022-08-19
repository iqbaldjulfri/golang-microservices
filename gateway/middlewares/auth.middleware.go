package middlewares

import (
	"gateway/interfaces"
	"gateway/models"
	"gateway/services/security"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type roleChecker interface {
	interfaces.FiberJSONSender
	interfaces.FiberLocalsGetterSetter
	interfaces.FiberNextRunner
	interfaces.FiberStatusSetter
}

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(security.JWT_SECRET),
		ErrorHandler:   security.JwtError,
		SuccessHandler: security.JwtSuccess,
		SigningMethod:  "HS512",
	})
}

func CheckRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return checkRoles(c, roles)
	}
}

func checkRoles(c roleChecker, roles []string) error {
	user := c.Locals("user").(*models.UserSafeDto)
	if user.Roles != nil && len(user.Roles) > 0 {
		for _, r1 := range roles {
			for _, r2 := range user.Roles {
				if r1 == r2.Code {
					return c.Next()
				}
			}
		}
	}

	return utils.JSONStatus(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, nil)
}

package handlers

import (
	mw "gateway/middlewares"
	"gateway/models"
	"gateway/services"
	"gateway/services/security"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
)

func AssignUsersHandlers(r fiber.Router) {
	group := r.Group("/users")

	group.Post("/login", validateLogin(), login)
	group.Get("/me", mw.Protected(), getMyProfile)
	group.Post("/", mw.Protected(), validateCreateUser(), createUser)
	group.Get("/", mw.Protected(), findUsers)
	group.Patch("/:id", mw.Protected(), validateUpdateUser(), updateUser)
}

func validateLogin() fiber.Handler {
	return mw.ValidateBodyFnFactory(func() interface{} {
		return new(models.LoginDto)
	})
}

func login(c *fiber.Ctx) error {
	dto := c.Locals("body").(*models.LoginDto)

	return security.DoLogin(c, *dto)
}

func getMyProfile(c *fiber.Ctx) error {
	user := security.GetUserFromLocals(c)

	return utils.JSON(c, user)
}

func createUser(c *fiber.Ctx) error {
	dto := c.Locals("body").(*models.CreateUserDto)
	user, err := services.CreateUser(*dto)

	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err, nil)
	}

	return utils.JSON(c, models.ToUserSafeDto(*user))
}

func validateCreateUser() fiber.Handler {
	return mw.ValidateBodyFnFactory(func() interface{} {
		return new(models.CreateUserDto)
	})
}

func findUsers(c *fiber.Ctx) error {
	username := c.Query("username")

	users, err := services.FindUsers(username)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err, nil)
	}

	var dtos []models.UserSafeDto
	for _, u := range users {
		dtos = append(dtos, *models.ToUserSafeDto(u))
	}

	return utils.JSON(c, dtos)
}

func updateUser(c *fiber.Ctx) error {
	tmpId, _ := c.ParamsInt("id")
	id := uint(tmpId)
	dto := c.Locals("body").(*models.UpdateUserDto)

	user, err := services.UpdateUser(id, *dto)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err, nil)
	}

	return utils.JSON(c, models.ToUserSafeDto(*user))
}

func validateUpdateUser() fiber.Handler {
	return mw.ValidateBodyFnFactory(func() interface{} {
		return new(models.UpdateUserDto)
	})
}

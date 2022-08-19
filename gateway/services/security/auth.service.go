package security

import (
	"errors"
	"gateway/models"
	"gateway/services"
	"gateway/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const JWT_SECRET = "secret"

func DoLogin(c *fiber.Ctx, loginDto models.LoginDto) error {
	user, err := authenticateUser(loginDto.Username, loginDto.Password)
	if err != nil {
		return utils.JSONError(c, fiber.StatusUnauthorized, err, nil)
	}

	token, err := generateJWT(*user)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err, nil)
	}

	return utils.JSON(c, models.LoginResponseDto{AccessToken: *token})
}

func authenticateUser(username, password string) (*models.User, error) {
	user, err := services.GetUserByUsername(username)
	if err != nil || !user.IsActive {
		return nil, errors.New("Authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Authentication failed")
	}

	return user, nil
}

func generateJWT(user models.User) (*string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	tokenizer := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := tokenizer.SignedString([]byte(JWT_SECRET))
	if err != nil {
		log.Fatal(err.Error())
		return nil, errors.New("Failed to sign JWT")
	}

	return &token, nil
}

func JwtError(c *fiber.Ctx, err error) error {
	return utils.JSONError(c, fiber.StatusUnauthorized, err, nil)
}

func JwtSuccess(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	user, err := services.GetUserByUsername(username)

	if err != nil {
		return utils.JSONStatus(c, fiber.StatusUnauthorized, "User not found", fiber.Map{username: username})
	}

	c.Locals("user", models.ToUserSafeDto(*user))

	return c.Next()
}

func GetUserFromLocals(c *fiber.Ctx) *models.UserSafeDto {
	return c.Locals("user").(*models.UserSafeDto)
}

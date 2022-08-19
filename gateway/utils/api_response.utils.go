package utils

import (
	"gateway/interfaces"

	"github.com/gofiber/fiber/v2"
)

type DefaultResponseBody struct {
	Status  int         `json:"statusCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type fiberJSONStatusSender interface {
	interfaces.FiberJSONSender
	interfaces.FiberStatusSetter
}

func JSONError(c fiberJSONStatusSender, status int, err error, data interface{}) error {
	return JSONStatus(c, status, err.Error(), data)
}

func JSONStatus(c fiberJSONStatusSender, status int, msg string, data interface{}) error {
	body := DefaultResponseBody{
		Status:  status,
		Message: msg,
		Data:    data,
	}

	c.Status(status)
	return c.JSON(body)
}

func JSON(c fiberJSONStatusSender, data interface{}) error {
	return JSONMessage(c, "Success", data)
}

func JSONMessage(c fiberJSONStatusSender, msg string, data interface{}) error {
	return JSONStatus(c, fiber.StatusOK, msg, data)
}

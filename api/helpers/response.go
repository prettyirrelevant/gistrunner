package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type SuccessPayload struct {
	StatusCode int    `json:"code"`
	Data       any    `json:"data"`
	Message    string `json:"message"`
}

type ErrorPayload struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func SuccessResponse(c *fiber.Ctx, payload *SuccessPayload) error {
	return c.Status(payload.StatusCode).JSON(payload)
}

func ErrorResponse(c *fiber.Ctx, payload *ErrorPayload) error {
	c.Locals("hasErr", "yes")
	return c.Status(payload.StatusCode).JSON(payload)
}

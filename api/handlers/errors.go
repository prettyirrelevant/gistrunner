package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/prettyirrelevant/gistrunner/helpers"
)

func GenericErrorHandler(c *fiber.Ctx, err error, logger *zap.Logger) error {
	logger.Error("an exception occurred", zap.Any("requestID", c.Locals("requestid")), zap.String("url", c.OriginalURL()), zap.String("method", c.Method()), zap.Error(err))
	return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusInternalServerError, Message: "Oops! Something went wrong on our end"})
}

func NotFoundHandler(c *fiber.Ctx) error {
	return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("%s %s resource does not exist on this server", c.Method(), c.OriginalURL())})
}

func TooManyRequestsHandler(c *fiber.Ctx) error {
	return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusTooManyRequests, Message: "Too many requests"})
}

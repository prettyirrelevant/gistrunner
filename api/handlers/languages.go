package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/gistrunner/helpers"
	"github.com/prettyirrelevant/gistrunner/types/languages"
)

func GetSupportedLanguages(c *fiber.Ctx) error {
	return helpers.SuccessResponse(c, &helpers.SuccessPayload{Message: "Supported languages returned successfully", Data: languages.SupportedLanguages, StatusCode: http.StatusOK})
}

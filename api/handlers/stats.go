package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/prettyirrelevant/gistrunner/database"
	"github.com/prettyirrelevant/gistrunner/helpers"
)

func GetStats(c *fiber.Ctx, db *database.Queries, logger *zap.Logger) error {
	count, err := db.CountGists(c.Context())
	if err != nil {
		logger.Error("could not get gists stat", zap.Error(err), zap.Any("requestID", c.Locals("requestid")))
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusInternalServerError, Message: err.Error()})
	}

	return helpers.SuccessResponse(c, &helpers.SuccessPayload{StatusCode: http.StatusOK, Data: map[string]int64{"count": count}, Message: "Gists count returned successfully"})
}

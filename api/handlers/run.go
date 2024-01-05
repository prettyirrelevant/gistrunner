package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/prettyirrelevant/gistrunner/database"
	"github.com/prettyirrelevant/gistrunner/helpers"
	"github.com/prettyirrelevant/gistrunner/services"
	"github.com/prettyirrelevant/gistrunner/types/languages"
)

var (
	ErrRunningGist          = errors.New("gist could not be executed")
	ErrContentCannotBeEmpty = errors.New("code content cannot be empty")
	ErrLanguageNotSupported = errors.New("language is not currently supported")
	ErrDatabaseNotReached   = errors.New("something broke on our end. try again later")
)

type RunGithubGistRequest struct {
	Content  string                        `json:"content"`
	Language languages.ProgrammingLanguage `json:"language"`
}

func (r *RunGithubGistRequest) Validate() error {
	if _, ok := languages.SupportedLanguages[r.Language]; !ok {
		return ErrLanguageNotSupported
	}

	if r.Content == "" {
		return ErrContentCannotBeEmpty
	}

	r.Content = strings.Trim(r.Content, " ")
	return nil
}

func RunGithubGist(c *fiber.Ctx, db *database.Queries, coderunner *services.CodeRunner, logger *zap.Logger) error {
	var requestBody RunGithubGistRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusUnprocessableEntity, Message: err.Error()})
	}

	if err := requestBody.Validate(); err != nil {
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusUnprocessableEntity, Message: err.Error()})
	}

	contentHash := helpers.HashString(requestBody.Content)
	gist, err := db.GetGist(c.Context(), contentHash)
	if err == nil {
		return helpers.SuccessResponse(c, &helpers.SuccessPayload{StatusCode: http.StatusOK, Message: "Gist ran successfully", Data: gist})
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		logger.Error("error fetching gist from db", zap.Error(err), zap.String("language", string(requestBody.Language)), zap.String("gist", requestBody.Content), zap.Any("requestID", c.Locals("requestid")))
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusInternalServerError, Message: ErrDatabaseNotReached.Error()})
	}

	output, err := coderunner.Run(requestBody.Language, requestBody.Content)
	if err != nil {
		logger.Error("error during code execution", zap.Error(err), zap.String("language", string(requestBody.Language)), zap.String("gist", requestBody.Content), zap.Any("requestID", c.Locals("requestid")))
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusBadGateway, Message: ErrRunningGist.Error()})
	}

	newGist, err := db.CreateGist(c.Context(), database.CreateGistParams{ID: helpers.GenerateID("gist"), Hash: contentHash, Language: string(requestBody.Language), Result: output})
	if err != nil {
		logger.Error("error creating gist in db", zap.Error(err), zap.String("language", string(requestBody.Language)), zap.String("gist", requestBody.Content), zap.Any("requestID", c.Locals("requestid")))
		return helpers.ErrorResponse(c, &helpers.ErrorPayload{StatusCode: http.StatusInternalServerError, Message: ErrRunningGist.Error()})
	}

	return helpers.SuccessResponse(c, &helpers.SuccessPayload{StatusCode: http.StatusOK, Message: "Gist ran successfully", Data: newGist})
}

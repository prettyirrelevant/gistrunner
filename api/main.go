package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/redis"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/prettyirrelevant/gistrunner/config"
	"github.com/prettyirrelevant/gistrunner/database"
	"github.com/prettyirrelevant/gistrunner/handlers"
	"github.com/prettyirrelevant/gistrunner/logger"
	"github.com/prettyirrelevant/gistrunner/services"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	// config
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("unable to init config", zap.Error(err))
	}

	// database
	dbCtx := context.Background()
	conn, err := pgx.Connect(dbCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}
	db := database.New(conn)

	// services
	coderunner, err := services.NewCodeRunner(cfg.DockerEngineURL, logger)
	if err != nil {
		logger.Fatal("unable to init code runner service", zap.Error(err))
	}

	var runInPreforkMode bool
	if cfg.Environment == "production" {
		runInPreforkMode = true
	}

	app := fiber.New(fiber.Config{
		Prefork: runInPreforkMode,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return handlers.GenericErrorHandler(c, err, logger)
		},
	})

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())

	app.Use(cache.New(cache.Config{
		CacheControl: true,
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/api/languages" {
				return true
			}
			if c.QueryBool("ignoreCache") {
				return true
			}
			// do not cache error responses
			// check `helpers.ErrorResponse()`
			if c.Locals("hasErr") == "yes" {
				return true
			}
			return false
		},
		Storage: redis.New(redis.Config{
			URL: cfg.RedisURL,
		}),
		Expiration: 1 * time.Hour,
	}))
	app.Use(limiter.New(limiter.Config{
		Storage: redis.New(redis.Config{
			URL: cfg.RedisURL,
		}),
		LimitReached: handlers.TooManyRequestsHandler,
	}))
	app.Use(healthcheck.New())

	app.Use(recover.New())

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		return handlers.GetStats(c, db, logger)
	})
	app.Post("/api/run", func(c *fiber.Ctx) error {
		return handlers.RunGithubGist(c, db, coderunner, logger)
	})

	app.Get("/api/languages", handlers.GetSupportedLanguages)

	app.Use(handlers.NotFoundHandler)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			logger.Panic("could not start server", zap.Error(err))
		}

		logger.Info("server running", zap.Int("port", cfg.Port))
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	logger.Info("gracefully shutting down server")
	_ = app.Shutdown()

	logger.Info("running cleanup tasks")

	// Your cleanup tasks go here
	conn.Close(dbCtx)
	coderunner.Close()
	logger.Info("server was successful shutdown")
}

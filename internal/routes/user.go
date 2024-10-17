package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/shubhamvscode/twitter-go-server/internal/handlers"
	"github.com/shubhamvscode/twitter-go-server/internal/middleware"
)

func SetupUserRoutes(app *fiber.App) {
	user := app.Group("/user", logger.New(), middleware.Auth())
	user.Get("/my-profile", handlers.GetMyProfile)
	// user.Get("/:username", handlers.GetProfile)
}

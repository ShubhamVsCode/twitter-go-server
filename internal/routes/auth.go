package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shubhamvscode/twitter-go-server/internal/handlers"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth", logger.New())
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
}

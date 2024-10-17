package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/shubhamvscode/twitter-go-server/internal/database"
	"github.com/shubhamvscode/twitter-go-server/internal/models"
	"github.com/shubhamvscode/twitter-go-server/internal/routes"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	err := database.DB.AutoMigrate(&models.User{}, &models.Tweet{}, &models.Comment{}, &models.Like{})
	if err != nil {
		log.Fatal(err)
	}

	routes.SetupAuthRoutes(app)
	routes.SetupTweetRoutes(app)
	routes.SetupUserRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"msg": "Hello, World!"})
	})

	log.Fatal(app.Listen(":8080"))
}

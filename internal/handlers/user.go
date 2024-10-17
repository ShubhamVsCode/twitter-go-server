package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shubhamvscode/twitter-go-server/internal/models"
)

func GetMyProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	username := user.Username

	return c.JSON(fiber.Map{
		"username": username,
	})
}

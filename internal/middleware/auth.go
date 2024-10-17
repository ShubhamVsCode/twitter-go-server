package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shubhamvscode/twitter-go-server/internal/database"
	"github.com/shubhamvscode/twitter-go-server/internal/models"
	"github.com/shubhamvscode/twitter-go-server/internal/utils"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenCookie := c.Cookies("token")
		if tokenCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token cookie",
			})
		}

		tokenString := strings.TrimPrefix(tokenCookie, "Bearer ")
		if tokenString == tokenCookie {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		userId, err := utils.ParseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		var user models.User
		if err := database.DB.First(&user, userId).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		c.Locals("user", user)
		return c.Next()
	}
}

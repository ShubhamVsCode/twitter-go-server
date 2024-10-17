package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shubhamvscode/twitter-go-server/internal/database"
	"github.com/shubhamvscode/twitter-go-server/internal/models"
	"github.com/shubhamvscode/twitter-go-server/internal/utils"
)

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserDoesNotExist         = errors.New("user does not exist")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrUsernamePasswordRequired = errors.New("username and password are required")
	ErrInternalServerError      = errors.New("internal server error")
)

func Register(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrUsernamePasswordRequired.Error()},
		)
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrUsernamePasswordRequired.Error()},
		)
	}

	database.DB.AutoMigrate(&models.User{})

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrUserAlreadyExists.Error()},
		)
	}

	if err := user.SetPassword(user.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": ErrInternalServerError.Error()},
		)
	}

	if err := database.DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": ErrInternalServerError.Error()},
		)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": ErrInternalServerError.Error()},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "Bearer " + token,
		Expires: time.Now().Add(time.Hour * 24),
		// HTTPOnly: true,
	})
	return c.JSON(
		&fiber.Map{"message": "User created successfully", "token": token},
	)
}

func Login(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": err.Error()},
		)
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrUsernamePasswordRequired.Error()},
		)
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrUserDoesNotExist.Error()},
		)
	}

	if !existingUser.CheckPassword(user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": ErrInvalidPassword.Error()},
		)
	}

	token, err := utils.GenerateToken(existingUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": ErrInternalServerError.Error()},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "Bearer " + token,
		Expires: time.Now().Add(time.Hour * 24),
		// HTTPOnly: true,
	})
	return c.JSON(
		&fiber.Map{"message": "Login successful", "token": token},
	)
}

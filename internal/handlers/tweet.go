package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shubhamvscode/twitter-go-server/internal/database"
	"github.com/shubhamvscode/twitter-go-server/internal/models"
	"gorm.io/gorm"
)

func CreateTweet(c *fiber.Ctx) error {
	tweet := new(models.Tweet)

	if err := c.BodyParser(tweet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": "Invalid request payload"},
		)
	}

	// Validate the tweet content
	if tweet.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": "Tweet content cannot be empty"},
		)
	}

	user := c.Locals("user").(models.User)
	tweet.UserID = user.ID

	// Save the tweet to the database
	if err := database.DB.Create(tweet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": "Failed to save tweet"},
		)
	}

	// Respond with the created tweet
	return c.Status(fiber.StatusCreated).JSON(tweet)
}

func GetTweets(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	user := c.Locals("user").(models.User)
	offset := (page - 1) * pageSize
	var tweets []struct {
		models.Tweet
		Username string `json:"username"`
		Likes    int    `json:"likes"`
		IsLiked  bool   `json:"is_liked"`
	}

	if err := database.DB.Table("tweets").
		Select("tweets.*, users.username, COUNT(DISTINCT likes.id) as likes, EXISTS(SELECT 1 FROM likes WHERE likes.tweet_id = tweets.id AND likes.user_id = ? ) as is_liked", user.ID).
		Joins("left join users on tweets.user_id = users.id").
		Joins("left join likes on tweets.id = likes.tweet_id AND likes.user_id = ?", user.ID).
		Group("tweets.id, users.username").
		Order("tweets.created_at desc").
		Offset(offset).Limit(pageSize).
		Find(&tweets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"error":   "Failed to fetch tweets",
				"message": err.Error(),
			},
		)
	}

	return c.JSON(tweets)
}

func CreateComment(c *fiber.Ctx) error {
	tweetID := c.Params("tweetId")
	comment := new(models.Comment)

	if err := c.BodyParser(comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": "Invalid request payload"},
		)
	}

	comment.TweetID = uuid.MustParse(tweetID)
	user := c.Locals("user").(models.User)
	comment.UserID = user.ID

	if err := database.DB.Create(comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": "Failed to save comment"},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func LikeTweet(c *fiber.Ctx) error {
	tweetID := c.Params("tweetId")
	isLiked := c.Query("like") == "true"
	user := c.Locals("user").(models.User)

	like := models.Like{
		UserID:  user.ID,
		TweetId: uuid.MustParse(tweetID),
	}

	if isLiked {
		var existingLike models.Like
		result := database.DB.Unscoped().Where("user_id = ? AND tweet_id = ?", user.ID, tweetID).First(&existingLike)

		if result.Error == nil {
			// Like exists, update deleted_at to null
			if err := database.DB.Model(&existingLike).Update("deleted_at", nil).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Map{
						"error":   "Failed to update like",
						"message": err.Error(),
					},
				)
			}
		} else if result.Error == gorm.ErrRecordNotFound {
			// Like doesn't exist, create new one
			if err := database.DB.Create(&like).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Map{
						"error":   "Failed to save like",
						"message": err.Error(),
					},
				)
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(
				&fiber.Map{
					"error":   "Database error",
					"message": result.Error.Error(),
				},
			)
		}
	} else {
		// Soft delete the like
		if err := database.DB.Where("user_id = ? AND tweet_id = ?", user.ID, tweetID).Delete(&like).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				&fiber.Map{
					"error":   "Failed to remove like",
					"message": err.Error(),
				},
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(like)
}

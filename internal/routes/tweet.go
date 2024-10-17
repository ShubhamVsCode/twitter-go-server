package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shubhamvscode/twitter-go-server/internal/handlers"
	"github.com/shubhamvscode/twitter-go-server/internal/middleware"
)

func SetupTweetRoutes(app *fiber.App) {
	tweet := app.Group("/tweet", logger.New(), middleware.Auth())
	tweet.Get("/", handlers.GetTweets)
	tweet.Post("/", handlers.CreateTweet)
	tweet.Post("/like/:tweetId", handlers.LikeTweet)
	tweet.Post("/comment/:tweetId", handlers.CreateComment)
}

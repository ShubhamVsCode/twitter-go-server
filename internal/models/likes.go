package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_tweet" json:"user_id"`
	TweetId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_tweet" json:"tweet_id"`
}

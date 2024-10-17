package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content   string    `gorm:"text;not null" json:"content"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TweetID   uuid.UUID `gorm:"type:uuid;not null" json:"tweet_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

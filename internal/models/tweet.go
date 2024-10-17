package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tweet struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content string    `gorm:"text;not null" json:"content"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}

func (t *Tweet) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

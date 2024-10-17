package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "postgresql://db_owner:gZB8mDWCzf7o@ep-blue-meadow-a14umq7i.ap-southeast-1.aws.neon.tech/db?sslmode=require"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

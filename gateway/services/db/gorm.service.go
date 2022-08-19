package db

import (
	"gateway/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB

func InitDB() {
	var err error
	Conn, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Panic("Failed to connect database")
	}

	Conn.AutoMigrate(
		&models.User{},
		&models.Role{},
	)
}

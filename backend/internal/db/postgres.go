package db

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func InitPostgres() {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}

		instance = conn
		log.Println("Connected to database")
	})
}

func GetDB() *gorm.DB {
	return instance
}

func SetDB(db *gorm.DB) {
	instance = db
}

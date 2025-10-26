package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set in .env file")
	}

	var db *gorm.DB
	var err error

	// Попытки подключения с задержкой
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			if i == maxRetries-1 {
				log.Fatal("Failed to connect to database after multiple attempts: ", err)
			}

			waitTime := time.Duration(i+1) * 2 * time.Second
			log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...",
				i+1, maxRetries, err, waitTime)
			time.Sleep(waitTime)
			continue
		}
		break
	}

	// Устанавливаем схему по умолчанию для этого соединения
	db.Exec("SET search_path TO task_schema")

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Database connection established")
	return db
}

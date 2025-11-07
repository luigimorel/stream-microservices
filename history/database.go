package main

import (
	"fmt"
	"os"
	"time"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (db *gorm.DB, err error) {
	gormLogger := slogGorm.New(slogGorm.WithIgnoreTrace())

	dbUsername := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Nairobi",
		dbHost, dbUsername, password, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		fmt.Printf("Database connection failed: %s\n", err)
		return nil, err
	}

	fmt.Println("Database connected successfully🚀")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Failed to get database instance: %s\n", err)
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}

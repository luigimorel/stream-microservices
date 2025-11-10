package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := ConnectDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("database has been connected successfully")
	}

	err = db.AutoMigrate(&Video{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	http.HandleFunc("/health", HealthHandler())
	http.HandleFunc("/history", HistoryHandler(db))
	http.HandleFunc("/save", SaveUserHistory(db))

	fmt.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

package main

import (
	"context"
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

	client, err := ConnectMongoDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("database has been connected successfully")
	}
	defer client.Disconnect(context.TODO())

	http.HandleFunc("/health", HealthHandler())
	http.HandleFunc("/videos", VideoHandler(client))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

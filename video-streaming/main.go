package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading .env file: %v\n", err)
		panic(err)
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

	fmt.Println("Server is running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

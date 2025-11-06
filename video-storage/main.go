package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	client, err := ConnectMongoDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("database has been connected successfully")
	}

	defer client.Disconnect(context.TODO())

	http.HandleFunc("/videos", VideoHandler(client))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

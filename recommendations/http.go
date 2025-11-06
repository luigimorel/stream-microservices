package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func VideoHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var video Video
		err := json.NewDecoder(r.Body).Decode(&video)
		if err != nil {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := StoreVideo(client, video); err != nil {
			slog.Error("Failed to store video", "error", err)
			http.Error(w, "Failed to store video", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Video stored successfully"})
	}
}

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "systems are okay",
		})
	}
}

func HistoryHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid user_id format", http.StatusBadRequest)
			return
		}

		videos, err := GetUserHistory(db, uint(userID))
		if err != nil {
			slog.Error("Failed to retrieve user history", "error", err)
			http.Error(w, "Failed to retrieve user history", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data":    videos,
			"message": "user history retrieved successfully",
		})
	}
}

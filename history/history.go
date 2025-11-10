package main

import (
	"context"
	"encoding/json"
	"fmt"
	"history/services"
	"log/slog"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Video struct {
	ID          uint       `gorm:"primaryKey"                         json:"id"`
	UserID      uint       `gorm:"not null;index"                     json:"user_id"`
	Title       string     `gorm:"not null"                           json:"title"`
	Description string     `                                          json:"description"`
	URL         string     `gorm:"not null"                           json:"url"`
	WatchedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"watched_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"                     json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"                     json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index"                              json:"deleted_at,omitempty"`
}

func GetUserHistory(db *gorm.DB, userID uint) ([]Video, error) {
	var videos []Video
	result := db.Where("user_id = ?", userID).Order("watched_at DESC").Find(&videos)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve user history: %v", result.Error)
	}

	return videos, nil
}

func SaveUserHistory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var video Video

		err := json.NewDecoder(r.Body).Decode(&video)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if video.UserID == 0 {
			http.Error(w, "user ID cannot be zero", http.StatusBadRequest)
		}
		if video.Title == "" || video.URL == "" {
			http.Error(w, "video title and URL cannot be empty", http.StatusBadRequest)
		}
		saveVideo := Video{
			Title:       video.Title,
			Description: video.Description,
			URL:         video.URL,
			WatchedAt:   time.Now(),
			UserID:      video.UserID,
		}

		result := db.Create(&saveVideo)
		if result.Error != nil {
			http.Error(w, "failed to save user history", http.StatusInternalServerError)
		}

		if err := saveMessageToQueue("video_history_queue", video); err != nil {
			http.Error(w, "failed to save message to queue", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Video stored successfully"})
	}
}

func saveMessageToQueue(queueName string, video Video) error {
	ch, err := services.Connect()
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)

	}

	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonBody, err := json.Marshal(video)
	if err != nil {
		return fmt.Errorf("failed to marshal video to JSON: %v", err)
	}

	err = ch.PublishWithContext(ctx, "", queueName, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        jsonBody,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

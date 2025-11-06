package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	URL         string     `gorm:"not null" json:"url"`
	WatchedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"watched_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func GetUserHistory(db *gorm.DB, userID uint) ([]Video, error) {
	var videos []Video
	result := db.Where("user_id = ?", userID).Order("watched_at DESC").Find(&videos)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve user history: %v", result.Error)
	}

	return videos, nil
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Video struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"        json:"id"`
	Title       string             `bson:"title"                json:"title"`
	Description string             `bson:"description"          json:"description"`
	URL         string             `bson:"url"                  json:"url"`
	CreatedAt   time.Time          `bson:"created_at"           json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"           json:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func StoreVideo(client *mongo.Client, video Video) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	saveVideo := Video{
		Id:          primitive.NewObjectID(),
		Title:       video.Title,
		Description: video.Description,
		URL:         video.URL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := client.Database(string(os.Getenv("DATABASE"))).Collection("videos")
	_, err := collection.InsertOne(ctx, saveVideo)
	if err != nil {
		return fmt.Errorf("failed to store video: %v", err)
	}
	return nil
}

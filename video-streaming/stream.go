package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"vid-stream/services"

	amqp "github.com/rabbitmq/amqp091-go"
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

	if err := publishVideo(saveVideo); err != nil {
		return fmt.Errorf("failed to publish video event: %v", err)
	}
	return nil
}

func publishVideo(video Video) error {
	ch, err := services.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	queueName := "video_created"

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	body, err := json.Marshal(video)
	if err != nil {
		return fmt.Errorf("failed to marshal video: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

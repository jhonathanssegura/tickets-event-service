package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-events/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-events/internal/db"
	"github.com/jhonathanssegura/ticket-events/internal/handler"
	"github.com/jhonathanssegura/ticket-events/internal/queue"
)

func main() {
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	queueURL := "http://localhost:4566/000000000000/event-queue"

	sqsClient := &queue.SQSClient{
		Client:   sqs.NewFromConfig(cfg),
		QueueURL: queueURL,
	}

	dynamoClient := &db.DynamoClient{
		Client: dynamodb.NewFromConfig(cfg),
	}

	handlerEvent := handler.NewEventHandler(sqsClient, dynamoClient)
	handlerCategory := handler.NewCategoryHandler(dynamoClient)
	// handlerQR := handler.NewQRHandler(dynamoClient)

	r := gin.Default()

	api := r.Group("/api")
	{
		// Event management endpoints
		api.GET("/events", handlerEvent.ListEvents)
		api.GET("/events/:id", handlerEvent.GetEvent)
		api.POST("/events", handlerEvent.CreateEvent)
		api.PUT("/events/:id", handlerEvent.UpdateEvent)
		api.DELETE("/events/:id", handlerEvent.DeleteEvent)
		// Category endpoint
		api.POST("/categories", handlerCategory.CreateCategory)
		// QR code endpoints eliminados
	}

	log.Println("🚀 Iniciando servidor de eventos en puerto 8080...")
	r.Run(":8080")
}

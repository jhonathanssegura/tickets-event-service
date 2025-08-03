package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-events/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-events/internal/db"
	"github.com/jhonathanssegura/ticket-events/internal/handler"
	"github.com/jhonathanssegura/ticket-events/internal/queue"
	"github.com/jhonathanssegura/ticket-events/internal/storage"
)

func main() {
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	queueURL := "http://localhost:4566/000000000000/event-queue"
	bucketName := "event-bucket"

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })

	log.Println("Verificando bucket S3...")
	_, err = s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Bucket S3 '%s' no existe, creándolo...", bucketName)
		_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			log.Fatalf("Error creando bucket S3: %v", err)
		}
		log.Printf("Bucket S3 '%s' creado exitosamente", bucketName)
	} else {
		log.Printf("Bucket S3 '%s' ya existe", bucketName)
	}

	sqsClient := &queue.SQSClient{
		Client:   sqs.NewFromConfig(cfg),
		QueueURL: queueURL,
	}

	storageClient := &storage.S3Client{
		Client:     s3Client,
		BucketName: bucketName,
	}

	dynamoClient := &db.DynamoClient{
		Client: dynamodb.NewFromConfig(cfg),
	}

	handlerEvent := handler.NewEventHandler(sqsClient, storageClient, dynamoClient)
	handlerCategory := handler.NewCategoryHandler(dynamoClient)
	handlerQR := handler.NewQRHandler(dynamoClient, storageClient)

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
		// QR code endpoints
		api.GET("/events/:id/qr", handlerQR.GetEventQR)
		api.GET("/events/:id/qr-s3", handlerQR.GetEventQRFromS3)
		api.POST("/qr/validate", handlerQR.ValidateQR)
		api.POST("/events/:id/qr", handlerQR.GenerateQRForEvent)
	}

	log.Println("🚀 Iniciando servidor de eventos en puerto 8080...")
	r.Run(":8080")
}

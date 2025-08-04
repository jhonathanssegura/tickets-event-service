package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-events/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-events/internal/model"
)

func main() {
	// Cargar configuración AWS
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	// Crear cliente DynamoDB
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Generar UUIDs para categorías
	categoryIDs := map[string]uuid.UUID{
		"cat-musica":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		"cat-teatro":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		"cat-deportes":   uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
		"cat-cine":       uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
		"cat-tecnologia": uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
	}

	// Datos de prueba para categorías
	categories := []model.Category{
		{
			ID:          categoryIDs["cat-musica"],
			Name:        "Música",
			Description: "Eventos musicales, conciertos y festivales",
			CreatedAt:   time.Now().Add(-7 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			ID:          categoryIDs["cat-teatro"],
			Name:        "Teatro",
			Description: "Obras de teatro, musicales y presentaciones escénicas",
			CreatedAt:   time.Now().Add(-6 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-6 * 24 * time.Hour),
		},
		{
			ID:          categoryIDs["cat-deportes"],
			Name:        "Deportes",
			Description: "Eventos deportivos, partidos y competiciones",
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:          categoryIDs["cat-cine"],
			Name:        "Cine",
			Description: "Estrenos de películas, festivales de cine y proyecciones especiales",
			CreatedAt:   time.Now().Add(-4 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-4 * 24 * time.Hour),
		},
		{
			ID:          categoryIDs["cat-tecnologia"],
			Name:        "Tecnología",
			Description: "Conferencias tecnológicas, hackathons y eventos de innovación",
			CreatedAt:   time.Now().Add(-3 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-3 * 24 * time.Hour),
		},
	}

	// Datos de prueba para eventos
	events := []model.Event{
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440101"),
			Name:        "Concierto de Rock en el Parque",
			Description: "Un increíble concierto de rock al aire libre con las mejores bandas del momento",
			CategoryID:  categoryIDs["cat-musica"],
			Location:    "Parque Central",
			Date:        time.Now().AddDate(0, 1, 15), // 1 mes y 15 días
			Capacity:    5000,
			Price:       75.00,
			Status:      model.EventStatusPublished,
			ImageURL:    "https://example.com/images/rock-concert.jpg",
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440102"),
			Name:        "Hamlet - Obra de Teatro Clásica",
			Description: "La famosa obra de Shakespeare presentada por la compañía nacional de teatro",
			CategoryID:  categoryIDs["cat-teatro"],
			Location:    "Teatro Nacional",
			Date:        time.Now().AddDate(0, 0, 10), // 10 días
			Capacity:    800,
			Price:       45.00,
			Status:      model.EventStatusPublished,
			ImageURL:    "https://example.com/images/hamlet.jpg",
			CreatedAt:   time.Now().Add(-20 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-20 * 24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440103"),
			Name:        "Final de Liga - Fútbol",
			Description: "La gran final de la liga local entre los dos mejores equipos",
			CategoryID:  categoryIDs["cat-deportes"],
			Location:    "Estadio Municipal",
			Date:        time.Now().AddDate(0, 0, 5), // 5 días
			Capacity:    25000,
			Price:       30.00,
			Status:      model.EventStatusPublished,
			ImageURL:    "https://example.com/images/football-final.jpg",
			CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-15 * 24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440104"),
			Name:        "Estreno Mundial - Nueva Película",
			Description: "El estreno mundial de la nueva película de acción y aventura",
			CategoryID:  categoryIDs["cat-cine"],
			Location:    "Cine Multiplex",
			Date:        time.Now().AddDate(0, 0, 3), // 3 días
			Capacity:    300,
			Price:       12.00,
			Status:      model.EventStatusPublished,
			ImageURL:    "https://example.com/images/movie-premiere.jpg",
			CreatedAt:   time.Now().Add(-10 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-10 * 24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440105"),
			Name:        "Conferencia de Tecnología 2024",
			Description: "La conferencia más importante del año sobre las últimas tendencias en tecnología",
			CategoryID:  categoryIDs["cat-tecnologia"],
			Location:    "Centro de Convenciones",
			Date:        time.Now().AddDate(0, 2, 0), // 2 meses
			Capacity:    1000,
			Price:       150.00,
			Status:      model.EventStatusPublished,
			ImageURL:    "https://example.com/images/tech-conference.jpg",
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
		},
	}

	fmt.Println("🌱 Cargando datos de prueba en DynamoDB...")

	// Insertar categorías en DynamoDB
	fmt.Printf("📊 Insertando %d categorías...\n", len(categories))
	for i, category := range categories {
		item := map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: category.ID.String()},
			"name":        &types.AttributeValueMemberS{Value: category.Name},
			"description": &types.AttributeValueMemberS{Value: category.Description},
			"created_at":  &types.AttributeValueMemberS{Value: category.CreatedAt.Format(time.RFC3339)},
			"updated_at":  &types.AttributeValueMemberS{Value: category.UpdatedAt.Format(time.RFC3339)},
		}

		_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("categories"),
			Item:      item,
		})

		if err != nil {
			log.Printf("Error insertando categoría %d: %v", i+1, err)
		} else {
			fmt.Printf("✅ Categoría '%s' insertada correctamente\n", category.Name)
		}
	}

	// Insertar eventos en DynamoDB
	fmt.Printf("\n📊 Insertando %d eventos...\n", len(events))
	for i, event := range events {
		item := map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: event.ID.String()},
			"name":        &types.AttributeValueMemberS{Value: event.Name},
			"description": &types.AttributeValueMemberS{Value: event.Description},
			"category_id": &types.AttributeValueMemberS{Value: event.CategoryID.String()},
			"location":    &types.AttributeValueMemberS{Value: event.Location},
			"date":        &types.AttributeValueMemberS{Value: event.Date.Format(time.RFC3339)},
			"capacity":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.Capacity)},
			"price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", event.Price)},
			"status":      &types.AttributeValueMemberS{Value: event.Status},
			"image_url":   &types.AttributeValueMemberS{Value: event.ImageURL},
			"created_at":  &types.AttributeValueMemberS{Value: event.CreatedAt.Format(time.RFC3339)},
			"updated_at":  &types.AttributeValueMemberS{Value: event.UpdatedAt.Format(time.RFC3339)},
		}

		_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("events"),
			Item:      item,
		})

		if err != nil {
			log.Printf("Error insertando evento %d: %v", i+1, err)
		} else {
			fmt.Printf("✅ Evento '%s' insertado correctamente\n", event.Name)
		}
	}

	fmt.Println("\n🎉 Datos de prueba cargados exitosamente!")
	fmt.Println("\n📋 Resumen de datos cargados:")
	fmt.Println("   • 5 categorías de eventos")
	fmt.Println("   • 5 eventos de prueba:")
	fmt.Println("     - Concierto de Rock en el Parque - $75.00")
	fmt.Println("     - Hamlet - Obra de Teatro Clásica - $45.00")
	fmt.Println("     - Final de Liga - Fútbol - $30.00")
	fmt.Println("     - Estreno Mundial - Nueva Película - $12.00")
	fmt.Println("     - Conferencia de Tecnología 2024 - $150.00")

	fmt.Println("\n🧪 Pruebas que puedes realizar:")
	fmt.Println("1. Listar todos los eventos:")
	fmt.Println("   curl -X GET http://localhost:8080/api/events")
	fmt.Println("\n2. Filtrar por categoría:")
	fmt.Println("   curl -X GET 'http://localhost:8080/api/events?category_id=550e8400-e29b-41d4-a716-446655440001'")
	fmt.Println("\n3. Obtener un evento específico:")
	fmt.Println("   curl -X GET http://localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440101")
	fmt.Println("\n4. Crear una nueva categoría:")
	fmt.Println("   curl -X POST http://localhost:8080/api/categories \\")
	fmt.Println("     -H 'Content-Type: application/json' \\")
	fmt.Println("     -d '{\"name\":\"Arte\",\"description\":\"Eventos de arte y exposiciones\"}'")
	fmt.Println("\n5. Crear un nuevo evento:")
	fmt.Println("   curl -X POST http://localhost:8080/api/events \\")
	fmt.Println("     -H 'Content-Type: application/json' \\")
	fmt.Println("     -d '{\"name\":\"Nuevo Evento\",\"description\":\"Descripción del evento\",\"category_id\":\"550e8400-e29b-41d4-a716-446655440001\",\"location\":\"Ubicación\",\"date\":\"2024-08-15T19:00:00Z\",\"capacity\":100,\"price\":25.00}'")
} 
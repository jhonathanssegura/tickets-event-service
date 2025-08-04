package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-events/internal/model"
)

type DynamoClient struct {
	Client *dynamodb.Client
}

func (d *DynamoClient) SaveEvent(event model.Event) error {
	fmt.Printf("Guardando evento: ID=%s, Name=%s, CategoryID=%s\n",
		event.ID.String(), event.Name, event.CategoryID.String())

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: event.ID.String()},
		"name":        &types.AttributeValueMemberS{Value: event.Name},
		"description": &types.AttributeValueMemberS{Value: event.Description},
		"category_id": &types.AttributeValueMemberS{Value: event.CategoryID.String()},
		"location":    &types.AttributeValueMemberS{Value: event.Location},
		"date":        &types.AttributeValueMemberS{Value: event.Date.Format(time.RFC3339)},
		"capacity":    &types.AttributeValueMemberN{Value: strconv.Itoa(event.Capacity)},
		"price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", event.Price)},
		"status":      &types.AttributeValueMemberS{Value: event.Status},
		"image_url":   &types.AttributeValueMemberS{Value: event.ImageURL},
		"created_at":  &types.AttributeValueMemberS{Value: event.CreatedAt.Format(time.RFC3339)},
		"updated_at":  &types.AttributeValueMemberS{Value: event.UpdatedAt.Format(time.RFC3339)},
	}

	_, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("events"),
		Item:      item,
	})

	if err != nil {
		var errorMsg string
		switch {
		case strings.Contains(err.Error(), "ResourceNotFoundException"):
			errorMsg = "La tabla 'events' no existe en DynamoDB. Verifique que LocalStack esté ejecutándose y la tabla haya sido creada."
		case strings.Contains(err.Error(), "RequestCanceled"):
			errorMsg = "Error de conexión con DynamoDB. Verifique que LocalStack esté ejecutándose en http://localhost:4566."
		case strings.Contains(err.Error(), "ConditionalCheckFailedException"):
			errorMsg = "El evento ya existe en la base de datos."
		default:
			errorMsg = fmt.Sprintf("Error guardando evento en DynamoDB: %v", err)
		}
		return fmt.Errorf(errorMsg)
	}

	return nil
}

func (d *DynamoClient) GetEventByID(eventID string) (*model.Event, error) {
	result, err := d.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("events"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("event not found")
	}

	event, err := d.unmarshalEvent(result.Item)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (d *DynamoClient) GetEvents(categoryID string, limit int) ([]model.Event, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String("events"),
		Limit:     aws.Int32(int32(limit)),
	}

	if categoryID != "" {
		scanInput.FilterExpression = aws.String("#category_id = :category_id")
		scanInput.ExpressionAttributeNames = map[string]string{
			"#category_id": "category_id",
		}
		categoryUUID, err := uuid.Parse(categoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID format: %v", err)
		}
		scanInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			":category_id": &types.AttributeValueMemberS{Value: categoryUUID.String()},
		}
	}

	result, err := d.Client.Scan(context.TODO(), scanInput)
	if err != nil {
		return nil, err
	}

	var events []model.Event
	for _, item := range result.Items {
		event, err := d.unmarshalEvent(item)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (d *DynamoClient) DeleteEvent(eventID string) error {
	_, err := d.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("events"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	return err
}

func (d *DynamoClient) SaveCategory(category model.Category) error {
	fmt.Printf("Guardando categoría: ID=%s, Name=%s\n", category.ID.String(), category.Name)

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: category.ID.String()},
		"name":        &types.AttributeValueMemberS{Value: category.Name},
		"description": &types.AttributeValueMemberS{Value: category.Description},
		"created_at":  &types.AttributeValueMemberS{Value: category.CreatedAt.Format(time.RFC3339)},
		"updated_at":  &types.AttributeValueMemberS{Value: category.UpdatedAt.Format(time.RFC3339)},
	}

	_, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("categories"),
		Item:      item,
	})

	if err != nil {
		var errorMsg string
		switch {
		case strings.Contains(err.Error(), "ResourceNotFoundException"):
			errorMsg = "La tabla 'categories' no existe en DynamoDB. Verifique que LocalStack esté ejecutándose y la tabla haya sido creada."
		case strings.Contains(err.Error(), "RequestCanceled"):
			errorMsg = "Error de conexión con DynamoDB. Verifique que LocalStack esté ejecutándose en http://localhost:4566."
		case strings.Contains(err.Error(), "ConditionalCheckFailedException"):
			errorMsg = "La categoría ya existe en la base de datos."
		default:
			errorMsg = fmt.Sprintf("Error guardando categoría en DynamoDB: %v", err)
		}
		return fmt.Errorf(errorMsg)
	}

	return nil
}

func (d *DynamoClient) unmarshalEvent(item map[string]types.AttributeValue) (*model.Event, error) {
	event := &model.Event{}

	if idVal, ok := item["id"].(*types.AttributeValueMemberS); ok {
		id, err := uuid.Parse(idVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid event ID: %v", err)
		}
		event.ID = id
	}

	if nameVal, ok := item["name"].(*types.AttributeValueMemberS); ok {
		event.Name = nameVal.Value
	}

	if descriptionVal, ok := item["description"].(*types.AttributeValueMemberS); ok {
		event.Description = descriptionVal.Value
	}

	if categoryIDVal, ok := item["category_id"].(*types.AttributeValueMemberS); ok {
		categoryID, err := uuid.Parse(categoryIDVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %v", err)
		}
		event.CategoryID = categoryID
	}

	if locationVal, ok := item["location"].(*types.AttributeValueMemberS); ok {
		event.Location = locationVal.Value
	}

	if dateVal, ok := item["date"].(*types.AttributeValueMemberS); ok {
		date, err := time.Parse(time.RFC3339, dateVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid date: %v", err)
		}
		event.Date = date
	}

	if capacityVal, ok := item["capacity"].(*types.AttributeValueMemberN); ok {
		capacity, err := strconv.Atoi(capacityVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid capacity: %v", err)
		}
		event.Capacity = capacity
	}

	if priceVal, ok := item["price"].(*types.AttributeValueMemberN); ok {
		price, err := strconv.ParseFloat(priceVal.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price: %v", err)
		}
		event.Price = price
	}

	if statusVal, ok := item["status"].(*types.AttributeValueMemberS); ok {
		event.Status = statusVal.Value
	}

	if imageURLVal, ok := item["image_url"].(*types.AttributeValueMemberS); ok {
		event.ImageURL = imageURLVal.Value
	}

	if createdAtVal, ok := item["created_at"].(*types.AttributeValueMemberS); ok {
		createdAt, err := time.Parse(time.RFC3339, createdAtVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid created_at time: %v", err)
		}
		event.CreatedAt = createdAt
	}

	if updatedAtVal, ok := item["updated_at"].(*types.AttributeValueMemberS); ok {
		updatedAt, err := time.Parse(time.RFC3339, updatedAtVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid updated_at time: %v", err)
		}
		event.UpdatedAt = updatedAt
	}

	return event, nil
} 
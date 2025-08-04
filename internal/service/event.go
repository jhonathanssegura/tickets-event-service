package service

import (
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-events/internal/db"
	"github.com/jhonathanssegura/ticket-events/internal/model"
)

type EventService struct {
	dynamoDB *db.DynamoClient
}

func NewEventService(dynamoDB *db.DynamoClient) *EventService {
	return &EventService{dynamoDB: dynamoDB}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(req model.CreateEventRequest) (*model.Event, error) {
	eventID := uuid.New()
	event := &model.Event{
		ID:          eventID,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Location:    req.Location,
		Date:        req.Date,
		Capacity:    req.Capacity,
		Price:       req.Price,
		Status:      model.EventStatusDraft,
		ImageURL:    req.ImageURL,
	}

	err := s.dynamoDB.SaveEvent(*event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(id string) (*model.Event, error) {
	return s.dynamoDB.GetEventByID(id)
}

// ListEvents retrieves all events with optional filtering
func (s *EventService) ListEvents(categoryID string, limit int) ([]model.Event, error) {
	return s.dynamoDB.GetEvents(categoryID, limit)
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(id string, updatedEvent *model.Event) (*model.Event, error) {
	// Get existing event
	existingEvent, err := s.dynamoDB.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if updatedEvent.Name != "" {
		existingEvent.Name = updatedEvent.Name
	}
	if updatedEvent.Description != "" {
		existingEvent.Description = updatedEvent.Description
	}
	if updatedEvent.CategoryID != uuid.Nil {
		existingEvent.CategoryID = updatedEvent.CategoryID
	}
	if updatedEvent.Location != "" {
		existingEvent.Location = updatedEvent.Location
	}
	if !updatedEvent.Date.IsZero() {
		existingEvent.Date = updatedEvent.Date
	}
	if updatedEvent.Capacity > 0 {
		existingEvent.Capacity = updatedEvent.Capacity
	}
	if updatedEvent.Price > 0 {
		existingEvent.Price = updatedEvent.Price
	}
	if updatedEvent.Status != "" {
		existingEvent.Status = updatedEvent.Status
	}
	if updatedEvent.ImageURL != "" {
		existingEvent.ImageURL = updatedEvent.ImageURL
	}

	err = s.dynamoDB.SaveEvent(*existingEvent)
	if err != nil {
		return nil, err
	}

	return existingEvent, nil
}

// DeleteEvent deletes an event
func (s *EventService) DeleteEvent(id string) error {
	return s.dynamoDB.DeleteEvent(id)
}

// GetEventWithStats retrieves an event with basic statistics
func (s *EventService) GetEventWithStats(id string) (*model.Event, error) {
	event, err := s.GetEvent(id)
	if err != nil {
		return nil, err
	}

	// For now, return the event as is since ticket statistics
	// would require integration with a ticket booking service
	return event, nil
}

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-events/internal/db"
	"github.com/jhonathanssegura/ticket-events/internal/model"
	"github.com/jhonathanssegura/ticket-events/internal/queue"
)

type EventHandler struct {
	DB  *db.DynamoClient
	SQS *queue.SQSClient
}

func NewEventHandler(sqs *queue.SQSClient, db *db.DynamoClient) *EventHandler {
	return &EventHandler{DB: db, SQS: sqs}
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	categoryID := c.Query("category_id")
	limitStr := c.Query("limit")
	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	events, err := h.DB.GetEvents(categoryID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo eventos", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
		"limit":  limit,
	})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento requerido"})
		return
	}

	event, err := h.DB.GetEventByID(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo evento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req model.CreateEventRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de evento inválidos", "details": err.Error()})
		return
	}

	eventID := uuid.New()
	now := time.Now()

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
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.DB.SaveEvent(*event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando evento", "details": err.Error()})
		return
	}

	// Send notification to SQS
	message := fmt.Sprintf("Nuevo evento creado: %s", event.Name)
	h.SQS.SendMessage(message)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Evento creado con éxito",
		"event":   event,
	})
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento requerido"})
		return
	}

	existingEvent, err := h.DB.GetEventByID(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo evento", "details": err.Error()})
		return
	}

	var req model.CreateEventRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización inválidos", "details": err.Error()})
		return
	}

	// Update fields
	if req.Name != "" {
		existingEvent.Name = req.Name
	}
	if req.Description != "" {
		existingEvent.Description = req.Description
	}
	if req.CategoryID != uuid.Nil {
		existingEvent.CategoryID = req.CategoryID
	}
	if req.Location != "" {
		existingEvent.Location = req.Location
	}
	if !req.Date.IsZero() {
		existingEvent.Date = req.Date
	}
	if req.Capacity > 0 {
		existingEvent.Capacity = req.Capacity
	}
	if req.Price > 0 {
		existingEvent.Price = req.Price
	}
	if req.ImageURL != "" {
		existingEvent.ImageURL = req.ImageURL
	}

	existingEvent.UpdatedAt = time.Now()

	if err := h.DB.SaveEvent(*existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando evento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Evento actualizado con éxito",
		"event":   existingEvent,
	})
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento requerido"})
		return
	}

	_, err := h.DB.GetEventByID(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando evento", "details": err.Error()})
		return
	}

	if err := h.DB.DeleteEvent(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando evento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evento eliminado con éxito"})
}

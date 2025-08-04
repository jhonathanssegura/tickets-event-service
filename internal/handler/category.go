package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-events/internal/db"
	"github.com/jhonathanssegura/ticket-events/internal/model"
)

type CategoryHandler struct {
	DB *db.DynamoClient
}

func NewCategoryHandler(db *db.DynamoClient) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de categoría inválidos", "details": err.Error()})
		return
	}

	categoryID := uuid.New()
	now := time.Now()

	category := &model.Category{
		ID:          categoryID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.DB.SaveCategory(*category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando categoría", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Categoría creada con éxito",
		"category": category,
	})
} 
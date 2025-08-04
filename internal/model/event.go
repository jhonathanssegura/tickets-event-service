package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CategoryID  uuid.UUID  `json:"category_id" db:"category_id"`
	Location    string     `json:"location" db:"location"`
	Date        time.Time  `json:"date" db:"date"`
	Capacity    int        `json:"capacity" db:"capacity"`
	Price       float64    `json:"price" db:"price"`
	Status      string     `json:"status" db:"status"`
	ImageURL    string     `json:"image_url" db:"image_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEventRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	ImageURL    string    `json:"image_url"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

const (
	EventStatusDraft     = "draft"
	EventStatusPublished = "published"
	EventStatusCancelled = "cancelled"
	EventStatusCompleted = "completed"
) 
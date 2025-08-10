package models

import (
	"time"
)

type FigmaFile struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	FileKey      string    `json:"file_key"`
	ImageURL     string    `json:"image_url"`
	Thumbnails   string    `json:"thumbnails,omitempty"`
	CanvasWidth  float64   `json:"canvas_width"`
	CanvasHeight float64   `json:"canvas_height"`
	ParsedAt     time.Time `json:"parsed_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Active       bool      `json:"active"`
}

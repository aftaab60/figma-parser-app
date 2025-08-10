package models

import (
	"encoding/json"
	"time"
)

// Component represents an extracted Figma component stored in the database.
// It corresponds to the 'components' table.
type Component struct {
	ID          int64           `json:"id"`
	FigmaFileID int64           `json:"figma_file_id"`
	NodeID      string          `json:"node_id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Description string          `json:"description,omitempty"`
	X           float64         `json:"x"`
	Y           float64         `json:"y"`
	Width       float64         `json:"width"`
	Height      float64         `json:"height"`
	ZIndex      int             `json:"z_index"`
	Properties  json.RawMessage `json:"properties,omitempty"` // Use json.RawMessage for JSONB
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Active      bool            `json:"active"`
}

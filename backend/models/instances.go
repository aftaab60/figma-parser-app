package models

import (
	"encoding/json"
	"time"
)

// Instance represents an instance of a component stored in the database.
// It corresponds to the 'instances' table.
type Instance struct {
	ID          int64           `json:"id"`
	ComponentID int64           `json:"component_id"`
	NodeID      string          `json:"node_id"`
	Name        string          `json:"name"`
	X           float64         `json:"x"`
	Y           float64         `json:"y"`
	Width       float64         `json:"width"`
	Height      float64         `json:"height"`
	Properties  json.RawMessage `json:"properties,omitempty"` // Use json.RawMessage for JSONB
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Active      bool            `json:"active"`
}

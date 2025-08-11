package figma_manager

import (
	"parser-service/models"
)

// ParsedFigmaData represents the complete parsed data from a Figma file
type ParsedFigmaData struct {
	File       *models.FigmaFile  `json:"file"`
	Components []models.Component `json:"components"`
	Instances  []models.Instance  `json:"instances"`
}

// FigmaAPIResponse represents the raw response from Figma API
type FigmaAPIResponse struct {
	Name          string                  `json:"name"`
	ThumbnailURL  string                  `json:"thumbnailUrl"`
	Document      Node                    `json:"document"`
	Components    map[string]Component    `json:"components"`
	ComponentSets map[string]ComponentSet `json:"componentSets"`
}

// Node represents a node in the Figma document tree
type Node struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name"`
	Type                string       `json:"type"`
	Visible             *bool        `json:"visible,omitempty"`
	ComponentID         string       `json:"componentId,omitempty"`
	AbsoluteBoundingBox *BoundingBox `json:"absoluteBoundingBox,omitempty"`
	Children            []Node       `json:"children,omitempty"`
	// Add other Figma properties as needed
}

// Component represents a Figma component from the API
type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Add other component properties
}

// ComponentSet represents a Figma component set from the API
type ComponentSet struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Add other component set properties
}

// BoundingBox represents position and dimensions
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

package figma_manager

import (
	"encoding/json"
	"fmt"
	"parser-service/models"
	"time"
)

// FigmaParser implements the IFigmaParser interface
type FigmaParser struct{}

// NewFigmaParser creates a new instance of FigmaParser
func NewFigmaParser() *FigmaParser {
	return &FigmaParser{}
}

// ParseFile parses a complete Figma API response into structured data
func (p *FigmaParser) ParseFile(apiResponse *FigmaAPIResponse, fileKey string, originalURL string) (*ParsedFigmaData, error) {
	if apiResponse == nil {
		return nil, fmt.Errorf("API response cannot be nil")
	}

	// Calculate canvas dimensions
	canvasWidth, canvasHeight, err := p.CalculateCanvasDimensions(apiResponse.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate canvas dimensions: %w", err)
	}

	// Create the FigmaFile model
	figmaFile := &models.FigmaFile{
		Name:         apiResponse.Name,
		URL:          originalURL, // Store the original URL
		FileKey:      fileKey,
		ImageURL:     apiResponse.ThumbnailURL,
		CanvasWidth:  canvasWidth,
		CanvasHeight: canvasHeight,
		ParsedAt:     time.Now(),
		Active:       true,
	}

	// Extract components from the API response
	components, err := p.extractComponentsFromAPI(apiResponse, 0) // fileID will be set later
	if err != nil {
		return nil, fmt.Errorf("failed to extract components: %w", err)
	}

	// Extract components from document nodes
	nodeComponents, err := p.ExtractComponents([]Node{apiResponse.Document}, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to extract components from nodes: %w", err)
	}

	// Combine components from API and nodes
	allComponents := append(components, nodeComponents...)

	// Deduplicate components by node_id to avoid constraint violations
	deduplicatedComponents := p.deduplicateComponents(allComponents)

	// Extract instances from document nodes
	instances, err := p.ExtractInstances([]Node{apiResponse.Document}, deduplicatedComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to extract instances: %w", err)
	}

	// Deduplicate instances by node_id to avoid constraint violations
	deduplicatedInstances := p.deduplicateInstances(instances)

	return &ParsedFigmaData{
		File:       figmaFile,
		Components: deduplicatedComponents,
		Instances:  deduplicatedInstances,
	}, nil
}

// ExtractComponents extracts components from Figma nodes
func (p *FigmaParser) ExtractComponents(nodes []Node, fileID int64) ([]models.Component, error) {
	var components []models.Component

	for _, node := range nodes {
		// Check if this node is a component
		if node.Type == "COMPONENT" || node.Type == "COMPONENT_SET" {
			component := p.nodeToComponent(node, fileID)
			components = append(components, component)
		}

		// Recursively process children
		if len(node.Children) > 0 {
			childComponents, err := p.ExtractComponents(node.Children, fileID)
			if err != nil {
				return nil, err
			}
			components = append(components, childComponents...)
		}
	}

	return components, nil
}

// ExtractInstances extracts instances from Figma nodes
func (p *FigmaParser) ExtractInstances(nodes []Node, components []models.Component) ([]models.Instance, error) {
	var instances []models.Instance

	for _, node := range nodes {
		// Check if this node is an instance
		if node.Type == "INSTANCE" && node.ComponentID != "" {
			// Find the corresponding component index for temporary linking
			componentIndex := p.findComponentIDByNodeID(node.ComponentID, components)
			if componentIndex > 0 {
				instance := p.nodeToInstance(node, componentIndex)
				instances = append(instances, instance)
			}
		}

		// Recursively process children
		if len(node.Children) > 0 {
			childInstances, err := p.ExtractInstances(node.Children, components)
			if err != nil {
				return nil, err
			}
			instances = append(instances, childInstances...)
		}
	}

	return instances, nil
}

// CalculateCanvasDimensions calculates the canvas dimensions from the document node
func (p *FigmaParser) CalculateCanvasDimensions(document Node) (width, height float64, err error) {
	if document.Type != "DOCUMENT" {
		return 0, 0, fmt.Errorf("root node must be of type DOCUMENT")
	}

	// Find the canvas node (usually the first child of document)
	for _, child := range document.Children {
		if child.Type == "CANVAS" && child.AbsoluteBoundingBox != nil {
			return child.AbsoluteBoundingBox.Width, child.AbsoluteBoundingBox.Height, nil
		}
	}

	// If no canvas found, calculate from all children
	var maxX, maxY float64
	p.calculateBounds(document.Children, &maxX, &maxY)

	if maxX == 0 && maxY == 0 {
		// Default dimensions if nothing found
		return 1920, 1080, nil
	}

	return maxX, maxY, nil
}

// extractComponentsFromAPI extracts components from the API's components map
func (p *FigmaParser) extractComponentsFromAPI(apiResponse *FigmaAPIResponse, fileID int64) ([]models.Component, error) {
	var components []models.Component

	// Extract from components map
	for nodeID, component := range apiResponse.Components {
		comp := models.Component{
			FigmaFileID: fileID,
			NodeID:      nodeID,
			Name:        component.Name,
			Type:        "COMPONENT",
			Description: component.Description,
			Active:      true,
		}

		// Find the actual node in the document tree to get position data
		if node := p.findNodeByID(apiResponse.Document, nodeID); node != nil {
			p.populateNodeProperties(node, &comp)
		}

		components = append(components, comp)
	}

	// Extract from component sets map
	for nodeID, componentSet := range apiResponse.ComponentSets {
		comp := models.Component{
			FigmaFileID: fileID,
			NodeID:      nodeID,
			Name:        componentSet.Name,
			Type:        "COMPONENT_SET",
			Description: componentSet.Description,
			Active:      true,
		}

		// Find the actual node in the document tree to get position data
		if node := p.findNodeByID(apiResponse.Document, nodeID); node != nil {
			p.populateNodeProperties(node, &comp)
		}

		components = append(components, comp)
	}

	return components, nil
}

// nodeToComponent converts a Figma node to a Component model
func (p *FigmaParser) nodeToComponent(node Node, fileID int64) models.Component {
	component := models.Component{
		FigmaFileID: fileID,
		NodeID:      node.ID,
		Name:        node.Name,
		Type:        node.Type,
		Active:      true,
	}

	p.populateNodeProperties(&node, &component)
	return component
}

// nodeToInstance converts a Figma node to an Instance model
// The componentIndex is temporary for parsing - actual DB relationships will be resolved during insertion
func (p *FigmaParser) nodeToInstance(node Node, componentIndex int64) models.Instance {
	instance := models.Instance{
		ComponentID: componentIndex, // Temporary index for parsing, will be resolved during DB insertion
		NodeID:      node.ID,
		Name:        node.Name,
		Active:      true,
	}

	// Set position and size
	if node.AbsoluteBoundingBox != nil {
		instance.X = node.AbsoluteBoundingBox.X
		instance.Y = node.AbsoluteBoundingBox.Y
		instance.Width = node.AbsoluteBoundingBox.Width
		instance.Height = node.AbsoluteBoundingBox.Height
	}

	// Set visibility
	if node.Visible != nil {
		instance.Active = *node.Visible
	}

	// Set z-index based on position in tree (simplified)
	instance.ZIndex = 0 // This could be enhanced with proper z-index calculation

	// Store the original Figma componentId in properties for later resolution
	additionalProps := make(map[string]interface{})
	additionalProps["figmaComponentId"] = node.ComponentID
	if node.Visible != nil {
		additionalProps["visible"] = *node.Visible
	}

	if len(additionalProps) > 0 {
		propsJSON, _ := json.Marshal(additionalProps)
		instance.Properties = propsJSON
	}

	return instance
}

// populateNodeProperties fills common properties from a Figma node
func (p *FigmaParser) populateNodeProperties(node *Node, component *models.Component) {
	// Set position and size
	if node.AbsoluteBoundingBox != nil {
		component.X = node.AbsoluteBoundingBox.X
		component.Y = node.AbsoluteBoundingBox.Y
		component.Width = node.AbsoluteBoundingBox.Width
		component.Height = node.AbsoluteBoundingBox.Height
	}

	// Set z-index based on position in tree (simplified)
	component.ZIndex = 0 // This could be enhanced with proper z-index calculation

	// Store any additional properties as JSON
	// This is where you could add more Figma-specific properties
	additionalProps := make(map[string]interface{})

	if node.Visible != nil {
		additionalProps["visible"] = *node.Visible
	}

	if node.Type != "" {
		additionalProps["nodeType"] = node.Type
	}

	if len(additionalProps) > 0 {
		propsJSON, _ := json.Marshal(additionalProps)
		component.Properties = propsJSON
	}
}

// findNodeByID recursively searches for a node with the given ID
func (p *FigmaParser) findNodeByID(node Node, targetID string) *Node {
	if node.ID == targetID {
		return &node
	}

	for _, child := range node.Children {
		if found := p.findNodeByID(child, targetID); found != nil {
			return found
		}
	}

	return nil
}

// findComponentIDByNodeID finds a component by its Figma node ID
// Returns the array index + 1 as a temporary reference for linking during parsing
func (p *FigmaParser) findComponentIDByNodeID(nodeID string, components []models.Component) int64 {
	for i, component := range components {
		if component.NodeID == nodeID {
			return int64(i + 1) // Temporary reference ID
		}
	}
	return 0
}

// calculateBounds recursively calculates the maximum bounds from nodes
func (p *FigmaParser) calculateBounds(nodes []Node, maxX, maxY *float64) {
	for _, node := range nodes {
		if node.AbsoluteBoundingBox != nil {
			nodeMaxX := node.AbsoluteBoundingBox.X + node.AbsoluteBoundingBox.Width
			nodeMaxY := node.AbsoluteBoundingBox.Y + node.AbsoluteBoundingBox.Height

			if nodeMaxX > *maxX {
				*maxX = nodeMaxX
			}
			if nodeMaxY > *maxY {
				*maxY = nodeMaxY
			}
		}

		// Recursively process children
		if len(node.Children) > 0 {
			p.calculateBounds(node.Children, maxX, maxY)
		}
	}
}

// assignZIndices assigns z-index values based on the order in the Figma tree
func (p *FigmaParser) assignZIndices(nodes []Node, startIndex int) int {
	currentIndex := startIndex

	for i := range nodes {
		// Assign z-index based on order (front to back in Figma)
		nodes[i] = p.setZIndexForNode(nodes[i], currentIndex)
		currentIndex++

		// Recursively assign for children
		if len(nodes[i].Children) > 0 {
			currentIndex = p.assignZIndices(nodes[i].Children, currentIndex)
		}
	}

	return currentIndex
}

// setZIndexForNode sets the z-index for a node (helper function)
func (p *FigmaParser) setZIndexForNode(node Node, zIndex int) Node {
	// This is a simplified approach - in a real implementation,
	// you might want to store this in the additional properties
	return node
}

// deduplicateComponents removes duplicate components based on node_id
// This prevents constraint violations when saving to database
func (p *FigmaParser) deduplicateComponents(components []models.Component) []models.Component {
	seen := make(map[string]bool)
	var deduplicated []models.Component

	for _, component := range components {
		if !seen[component.NodeID] {
			seen[component.NodeID] = true
			deduplicated = append(deduplicated, component)
		}
	}

	return deduplicated
}

// deduplicateInstances removes duplicate instances based on node_id
// This prevents constraint violations when saving to database
func (p *FigmaParser) deduplicateInstances(instances []models.Instance) []models.Instance {
	seen := make(map[string]bool)
	var deduplicated []models.Instance

	for _, instance := range instances {
		if !seen[instance.NodeID] {
			seen[instance.NodeID] = true
			deduplicated = append(deduplicated, instance)
		}
	}

	return deduplicated
}

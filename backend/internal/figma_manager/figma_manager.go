package figma_manager

import (
	"context"
	"fmt"
	"parser-service/models"
	"regexp"
	"strings"
)

// Having FigmaManager as iterface just in case we want to mock it in tests OR use a different implementation in the future
type IFigmaManager interface {
	ParseFigmaFileFromURL(ctx context.Context, figmaURL string) (*ParsedFigmaData, error)
	ParseFigmaFileFromKey(ctx context.Context, fileKey string) (*ParsedFigmaData, error)

	ExtractComponentsFromFile(ctx context.Context, fileKey string) ([]models.Component, error)
	ExtractInstancesFromFile(ctx context.Context, fileKey string) ([]models.Instance, error)
	ParseFigmaFileWithImages(ctx context.Context, fileKey string) (*ParsedFigmaData, map[string]string, error)
	GetFileImages(ctx context.Context, fileKey string, nodeIDs []string) (map[string]string, error)
	GetFileNodes(ctx context.Context, fileKey string, nodeIDs []string) (*FigmaAPIResponse, error)
	ValidateFileAccess(ctx context.Context, fileKey string) error

	ValidateFigmaToken(token string) error
}

// FigmaManager implements the IFigmaManager interface
type FigmaManager struct {
	client *FigmaClient
	parser *FigmaParser
}

func NewFigmaManager() *FigmaManager {
	client := NewFigmaClient()
	parser := NewFigmaParser()

	return &FigmaManager{
		client: client,
		parser: parser,
	}
}

// ParseFigmaFileFromURL parses a Figma file from URL and returns structured data
func (m *FigmaManager) ParseFigmaFileFromURL(ctx context.Context, figmaURL string) (*ParsedFigmaData, error) {
	// Extract file key from URL
	fileKey, err := m.extractFileKeyFromURL(figmaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract file key from URL: %w", err)
	}

	// Get file data from Figma API
	apiResponse, err := m.client.GetFile(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Figma API: %w", err)
	}

	// Parse the API response into our models, passing the original URL
	parsedData, err := m.parser.ParseFile(apiResponse, fileKey, figmaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file data: %w", err)
	}

	return parsedData, nil
}

// ParseFigmaFileFromKey parses a Figma file from file key and returns structured data
func (m *FigmaManager) ParseFigmaFileFromKey(ctx context.Context, fileKey string) (*ParsedFigmaData, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	// Get file data from Figma API
	apiResponse, err := m.client.GetFile(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Figma API: %w", err)
	}

	// Parse the API response into our models (no original URL available when parsing from key only)
	parsedData, err := m.parser.ParseFile(apiResponse, fileKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse file data: %w", err)
	}

	return parsedData, nil
}

// ExtractComponentsFromFile extracts only components from a Figma file
func (m *FigmaManager) ExtractComponentsFromFile(ctx context.Context, fileKey string) ([]models.Component, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	// Get file data from Figma API
	apiResponse, err := m.client.GetFile(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Figma API: %w", err)
	}

	// Extract components from the document
	components, err := m.parser.ExtractComponents([]Node{apiResponse.Document}, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to extract components: %w", err)
	}

	return components, nil
}

// ExtractInstancesFromFile extracts only instances from a Figma file
func (m *FigmaManager) ExtractInstancesFromFile(ctx context.Context, fileKey string) ([]models.Instance, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	// Get file data from Figma API
	apiResponse, err := m.client.GetFile(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Figma API: %w", err)
	}

	// First extract components to reference in instances
	components, err := m.parser.ExtractComponents([]Node{apiResponse.Document}, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to extract components: %w", err)
	}

	// Extract instances from the document
	instances, err := m.parser.ExtractInstances([]Node{apiResponse.Document}, components)
	if err != nil {
		return nil, fmt.Errorf("failed to extract instances: %w", err)
	}

	return instances, nil
}

// extractFileKeyFromURL extracts the file key from a Figma URL
func (m *FigmaManager) extractFileKeyFromURL(figmaURL string) (string, error) {
	if figmaURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Clean the URL
	figmaURL = strings.TrimSpace(figmaURL)

	// Figma file URL patterns:
	// https://www.figma.com/file/{file-key}/{file-name}
	// https://www.figma.com/design/{file-key}/{file-name}
	// https://figma.com/file/{file-key}/{file-name}

	patterns := []string{
		`(?:https?://)?(?:www\.)?figma\.com/file/([a-zA-Z0-9]+)`,
		`(?:https?://)?(?:www\.)?figma\.com/design/([a-zA-Z0-9]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(figmaURL)
		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	// If no pattern matches, check if it's already just a file key
	if isValidFileKey(figmaURL) {
		return figmaURL, nil
	}

	return "", fmt.Errorf("invalid Figma URL format: %s", figmaURL)
}

// isValidFileKey checks if a string looks like a valid Figma file key
func isValidFileKey(key string) bool {
	// Figma file keys are typically alphanumeric strings
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, key)
	return matched && len(key) > 10 // Figma keys are usually longer than 10 characters
}

// GetFileImages gets images for specific nodes (utility method)
func (m *FigmaManager) GetFileImages(ctx context.Context, fileKey string, nodeIDs []string) (map[string]string, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	return m.client.GetFileImages(ctx, fileKey, nodeIDs)
}

// GetFileNodes gets specific nodes from a file (utility method)
func (m *FigmaManager) GetFileNodes(ctx context.Context, fileKey string, nodeIDs []string) (*FigmaAPIResponse, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	return m.client.GetFileNodes(ctx, fileKey, nodeIDs)
}

// ValidateFileAccess checks if we can access a Figma file
func (m *FigmaManager) ValidateFileAccess(ctx context.Context, fileKey string) error {
	if fileKey == "" {
		return fmt.Errorf("file key cannot be empty")
	}

	// Try to get basic file info
	_, err := m.client.GetFile(ctx, fileKey)
	if err != nil {
		return fmt.Errorf("cannot access Figma file: %w", err)
	}

	return nil
}

// ParseFigmaFileWithImages parses a file and also fetches images for components
func (m *FigmaManager) ParseFigmaFileWithImages(ctx context.Context, fileKey string) (*ParsedFigmaData, map[string]string, error) {
	// First parse the file normally
	parsedData, err := m.ParseFigmaFileFromKey(ctx, fileKey)
	if err != nil {
		return nil, nil, err
	}

	// Collect all node IDs for image generation
	var nodeIDs []string

	// Add component node IDs
	for _, component := range parsedData.Components {
		if component.NodeID != "" {
			nodeIDs = append(nodeIDs, component.NodeID)
		}
	}

	// Add instance node IDs
	for _, instance := range parsedData.Instances {
		if instance.NodeID != "" {
			nodeIDs = append(nodeIDs, instance.NodeID)
		}
	}

	// Get images if we have node IDs
	var images map[string]string
	if len(nodeIDs) > 0 {
		images, err = m.client.GetFileImages(ctx, fileKey, nodeIDs)
		if err != nil {
			// Don't fail the whole operation if images fail
			// Just log and continue
			images = make(map[string]string)
		}
	} else {
		images = make(map[string]string)
	}

	return parsedData, images, nil
}

func (m *FigmaManager) ValidateFigmaToken(token string) error {
	if token == "" {
		return fmt.Errorf("Figma token cannot be empty")
	}

	// Use the client to validate the token
	return m.client.ValidateAPIToken(token)
}

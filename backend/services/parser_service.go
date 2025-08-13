package services

import (
	"context"
	"fmt"
	"parser-service/internal/figma_manager"
	"parser-service/models"
	"parser-service/repositories"
)

type ParserService struct {
	FigmaManager         figma_manager.IFigmaManager
	FigmaFilesRepository repositories.IFigmaFilesRepository
	ComponentsRepository repositories.IComponentsRepository
	InstancesRepository  repositories.IInstancesRepository
}

func NewParserService(
	figmaManager figma_manager.IFigmaManager,
	figmaFilesRepo repositories.IFigmaFilesRepository,
	componentsRepo repositories.IComponentsRepository,
	instancesRepo repositories.IInstancesRepository,
) *ParserService {
	return &ParserService{
		FigmaManager:         figmaManager,
		FigmaFilesRepository: figmaFilesRepo,
		ComponentsRepository: componentsRepo,
		InstancesRepository:  instancesRepo,
	}
}

// ParseAndSaveFigmaFile - Main method that accepts Figma URL and saves all extracted data
func (s *ParserService) ParseAndSaveFigmaFile(ctx context.Context, figmaURL string) (*models.FigmaFile, error) {
	// 1. Parse Figma file from URL
	parsedData, err := s.FigmaManager.ParseFigmaFileFromURL(ctx, figmaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Figma file: %w", err)
	}

	// 2. Save the Figma file record first
	savedFile, err := s.FigmaFilesRepository.CreateFigmaFile(ctx, parsedData.File)
	if err != nil {
		return nil, fmt.Errorf("failed to save Figma file: %w", err)
	}

	// 3. Save components with the file ID
	savedComponents, err := s.saveComponents(ctx, parsedData.Components, savedFile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to save components: %w", err)
	}

	// 4. Save instances with proper component references
	err = s.saveInstances(ctx, parsedData.Instances, savedComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to save instances: %w", err)
	}

	return savedFile, nil
}

// GetFigmaFileWithDetails - Retrieve a complete Figma file with components and instances
func (s *ParserService) GetFigmaFileWithDetails(ctx context.Context, fileID int64) (*FigmaFileDetails, error) {
	// Get file
	file, err := s.FigmaFilesRepository.GetFigmaFileByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Get components
	components, err := s.ComponentsRepository.GetComponentsByFigmaFileID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	// Get instances for the file
	instances, err := s.InstancesRepository.GetInstancesByFigmaFileID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	return &FigmaFileDetails{
		File:       file,
		Components: components,
		Instances:  instances,
	}, nil
}

// GetComponentsByFileID - Get all components for a specific file
func (s *ParserService) GetComponentsByFileID(ctx context.Context, fileID int64) ([]models.Component, error) {
	return s.ComponentsRepository.GetComponentsByFigmaFileID(ctx, fileID)
}

// GetInstancesByComponentID - Get all instances for a specific component
func (s *ParserService) GetInstancesByComponentID(ctx context.Context, componentID int64) ([]models.Instance, error) {
	return s.InstancesRepository.GetInstancesByComponentID(ctx, componentID)
}

// GetInstancesByFileID - Get all instances for a specific file
func (s *ParserService) GetInstancesByFileID(ctx context.Context, fileID int64) ([]models.Instance, error) {
	return s.InstancesRepository.GetInstancesByFigmaFileID(ctx, fileID)
}

// ValidateFigmaURL - Check if a Figma URL is accessible before parsing
func (s *ParserService) ValidateFigmaURL(ctx context.Context, figmaURL string) error {
	// Extract file key from URL and validate access
	parsedData, err := s.FigmaManager.ParseFigmaFileFromURL(ctx, figmaURL)
	if err != nil {
		return fmt.Errorf("invalid Figma URL or inaccessible file: %w", err)
	}

	// Just validate access, don't save anything
	return s.FigmaManager.ValidateFileAccess(ctx, parsedData.File.FileKey)
}

// saveComponents saves components and returns the saved components with database IDs
func (s *ParserService) saveComponents(ctx context.Context, components []models.Component, fileID int64) ([]models.Component, error) {
	var savedComponents []models.Component

	for _, component := range components {
		// Set the file ID for database foreign key
		component.FigmaFileID = fileID

		savedComponent, err := s.ComponentsRepository.CreateComponent(ctx, &component)
		if err != nil {
			return nil, fmt.Errorf("failed to save component %s: %w", component.Name, err)
		}

		savedComponents = append(savedComponents, *savedComponent)
	}

	return savedComponents, nil
}

// saveInstances saves instances and resolves component relationships
func (s *ParserService) saveInstances(ctx context.Context, instances []models.Instance, savedComponents []models.Component) error {
	// Create a map of temporary component IDs to actual database IDs
	componentIDMap := make(map[int64]int64)
	for i, component := range savedComponents {
		// The parser used (index + 1) as temporary ID
		tempID := int64(i + 1)
		componentIDMap[tempID] = component.ID
	}

	for _, instance := range instances {
		// Resolve the temporary component ID to actual database ID
		if actualComponentID, exists := componentIDMap[instance.ComponentID]; exists {
			instance.ComponentID = actualComponentID

			_, err := s.InstancesRepository.CreateInstance(ctx, &instance)
			if err != nil {
				return fmt.Errorf("failed to save instance %s: %w", instance.Name, err)
			}
		} else {
			// Log warning but don't fail - this might happen if component wasn't found during parsing
			fmt.Printf("Warning: Could not resolve component ID %d for instance %s\n", instance.ComponentID, instance.Name)
		}
	}

	return nil
}

// FigmaFileDetails represents a complete Figma file with all related data
type FigmaFileDetails struct {
	File       *models.FigmaFile  `json:"file"`
	Components []models.Component `json:"components"`
	Instances  []models.Instance  `json:"instances"`
}

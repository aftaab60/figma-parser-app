package figma_manager_test

import (
	"context"
	"os"
	"parser-service/internal/figma_manager"
	"testing"
)

// Test with sample Figma data
func TestFigmaManager_Integration_SampleData(t *testing.T) {
	apiToken := os.Getenv("FIGMA_API_TOKEN")
	if apiToken == "" {
		t.Fatal("FIGMA_API_TOKEN environment variable not set")
	}
	ctx := context.WithValue(context.Background(), "figma_token", apiToken)
	figmaURL := "https://www.figma.com/design/DNCLfE7Tf8A0mudOOLZUYx/Locofy.ai.test?t=94QjDVDsuj5cHJEU-1"

	t.Run("Test URL Extraction", func(t *testing.T) {
		manager := figma_manager.NewFigmaManager()

		// Test parsing through ParseFigmaFileFromURL to verify URL extraction
		parsedData, err := manager.ParseFigmaFileFromURL(ctx, figmaURL)
		if err != nil {
			t.Fatalf("Failed to parse Figma file from URL: %v", err)
		}

		if parsedData.File == nil {
			t.Fatal("Expected file data to be non-nil")
		}

		expectedFileKey := "DNCLfE7Tf8A0mudOOLZUYx"
		if parsedData.File.FileKey != expectedFileKey {
			t.Errorf("Expected file key %s, got %s", expectedFileKey, parsedData.File.FileKey)
		}

		t.Logf("Successfully extracted file key: %s", parsedData.File.FileKey)
	})

	t.Run("Test Complete File Parsing", func(t *testing.T) {
		manager := figma_manager.NewFigmaManager()

		parsedData, err := manager.ParseFigmaFileFromURL(ctx, figmaURL)
		if err != nil {
			t.Fatalf("Failed to parse Figma file: %v", err)
		}

		// Verify parsed data structure
		if parsedData.File == nil {
			t.Error("Expected file data to be non-nil")
		}

		if parsedData.File.Name == "" {
			t.Error("Expected file name to be non-empty")
		}

		if parsedData.File.CanvasWidth <= 0 || parsedData.File.CanvasHeight <= 0 {
			t.Error("Expected positive canvas dimensions")
		}

		t.Logf("Successfully parsed file: %s", parsedData.File.Name)
		t.Logf("Found %d components and %d instances", len(parsedData.Components), len(parsedData.Instances))

		// debugging - log 5 component details and rest with ...
		if len(parsedData.Components) > 0 {
			t.Logf("Components:")
			for i, component := range parsedData.Components[:min(5, len(parsedData.Components))] {
				t.Logf("%d. %s (%s) - Position: (%.1f, %.1f), Size: %.1fx%.1f",
					i+1, component.Name, component.Type,
					component.X, component.Y, component.Width, component.Height)
			}
			if len(parsedData.Components) > 5 {
				t.Logf("... and %d more components", len(parsedData.Components)-5)
			}
		}

		// Log instance details
		if len(parsedData.Instances) > 0 {
			t.Logf("Instances:")
			for i, instance := range parsedData.Instances[:min(5, len(parsedData.Instances))] {
				t.Logf("%d. %s - Position: (%.1f, %.1f), Size: %.1fx%.1f",
					i+1, instance.Name,
					instance.X, instance.Y, instance.Width, instance.Height)
			}
			if len(parsedData.Instances) > 5 {
				t.Logf("... and %d more instances", len(parsedData.Instances)-5)
			}
		}
	})

	t.Run("Test Components Only Extraction", func(t *testing.T) {
		manager := figma_manager.NewFigmaManager()
		fileKey := "DNCLfE7Tf8A0mudOOLZUYx"

		components, err := manager.ExtractComponentsFromFile(ctx, fileKey)
		if err != nil {
			t.Fatalf("Failed to extract components: %v", err)
		}

		t.Logf("Extracted %d components only", len(components))

		for _, component := range components {
			if component.NodeID == "" {
				t.Error("Component missing NodeID")
			}
			if component.Name == "" {
				t.Error("Component missing Name")
			}
			if component.Type == "" {
				t.Error("Component missing Type")
			}
		}
	})

	t.Run("Test Instances Only Extraction", func(t *testing.T) {
		manager := figma_manager.NewFigmaManager()
		fileKey := "DNCLfE7Tf8A0mudOOLZUYx"

		// First extract components to see what's available
		components, err := manager.ExtractComponentsFromFile(ctx, fileKey)
		if err != nil {
			t.Fatalf("Failed to extract components: %v", err)
		}
		t.Logf("Found %d components for instance matching", len(components))

		// Log component node IDs for debugging
		if len(components) > 0 {
			t.Logf("Component NodeIDs:")
			for i, comp := range components[:min(3, len(components))] {
				t.Logf("  %d. %s (NodeID: %s, Type: %s)", i+1, comp.Name, comp.NodeID, comp.Type)
			}
		}

		instances, err := manager.ExtractInstancesFromFile(ctx, fileKey)
		if err != nil {
			t.Fatalf("Failed to extract instances: %v", err)
		}

		t.Logf("Extracted %d instances only", len(instances))

		// Validate instance properties
		for _, instance := range instances {
			if instance.NodeID == "" {
				t.Error("Instance missing NodeID")
			}
			if instance.Name == "" {
				t.Error("Instance missing Name")
			}
			if instance.ComponentID <= 0 {
				t.Error("Instance missing or invalid ComponentID")
			}
			// Validate position and size
			if instance.Width <= 0 || instance.Height <= 0 {
				t.Errorf("Instance %s has invalid dimensions: %.1fx%.1f", instance.Name, instance.Width, instance.Height)
			}
		}

		// Log first few instances for debugging
		if len(instances) > 0 {
			t.Logf("First few instances:")
			for i, instance := range instances[:min(3, len(instances))] {
				t.Logf("%d. %s (Component ID: %d) - Position: (%.1f, %.1f), Size: %.1fx%.1f",
					i+1, instance.Name, instance.ComponentID,
					instance.X, instance.Y, instance.Width, instance.Height)
			}
			if len(instances) > 3 {
				t.Logf("... and %d more instances", len(instances)-3)
			}
		}
	})

	t.Run("Test File Parsing With Images", func(t *testing.T) {
		manager := figma_manager.NewFigmaManager()
		fileKey := "DNCLfE7Tf8A0mudOOLZUYx"

		parsedData, images, err := manager.ParseFigmaFileWithImages(ctx, fileKey)
		if err != nil {
			t.Fatalf("Failed to parse file with images: %v", err)
		}

		if parsedData == nil {
			t.Error("Expected parsed data to be non-nil")
		}

		t.Logf("Parsed file with %d images", len(images))

		// Log image URLs (first few)
		imageCount := 0
		for nodeID, imageURL := range images {
			if imageCount >= 3 { // Limit output
				break
			}
			t.Logf("Node %s: %s", nodeID, imageURL)
			imageCount++
		}
		if len(images) > 3 {
			t.Logf("... and %d more images", len(images)-3)
		}

		// Verify images have valid URLs
		for nodeID, imageURL := range images {
			if imageURL == "" {
				t.Errorf("Empty image URL for node %s", nodeID)
			}
		}
	})
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

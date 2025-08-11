package figma_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"parser-service/internal/errors"
	"strings"
	"time"
)

const (
	FigmaBaseURL   = "https://api.figma.com/v1"
	DefaultTimeout = 30 * time.Second
)

// FigmaClient implements the IFigmaClient interface
type FigmaClient struct {
	apiToken   string
	httpClient *http.Client
	baseURL    string
}

// NewFigmaClient creates a new Figma API client
func NewFigmaClient(apiToken string) *FigmaClient {
	return &FigmaClient{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: FigmaBaseURL,
	}
}

// NewFigmaClientWithTimeout creates a new Figma API client with custom timeout
func NewFigmaClientWithTimeout(apiToken string, timeout time.Duration) *FigmaClient {
	return &FigmaClient{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: FigmaBaseURL,
	}
}

// GetFile retrieves complete file data from Figma API
func (c *FigmaClient) GetFile(fileKeyOrURL string) (*FigmaAPIResponse, error) {
	if fileKeyOrURL == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}

	// Clean the file key (remove any URL parts if a full URL was provided)
	fileKey := c.extractFileKeyFromURL(fileKeyOrURL)
	endpoint := fmt.Sprintf("%s/files/%s", c.baseURL, fileKey)

	response, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Figma API: %w", err)
	}

	var figmaResponse FigmaAPIResponse
	if err := json.Unmarshal(response, &figmaResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Figma API response: %w", err)
	}

	return &figmaResponse, nil
}

// GetFileNodes retrieves specific nodes from a Figma file
func (c *FigmaClient) GetFileNodes(fileKey string, nodeIDs []string) (*FigmaAPIResponse, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}
	if len(nodeIDs) == 0 {
		return nil, fmt.Errorf("node IDs cannot be empty")
	}

	// Clean the file key
	fileKey = c.extractFileKeyFromURL(fileKey)

	// Build the endpoint with node IDs as query parameters
	endpoint := fmt.Sprintf("%s/files/%s/nodes", c.baseURL, fileKey)

	// Add node IDs as query parameters
	params := url.Values{}
	params.Add("ids", strings.Join(nodeIDs, ","))
	endpoint += "?" + params.Encode()

	response, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file nodes from Figma API: %w", err)
	}

	var figmaResponse FigmaAPIResponse
	if err := json.Unmarshal(response, &figmaResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Figma API response: %w", err)
	}

	return &figmaResponse, nil
}

// GetFileImages retrieves rendered images for specific nodes
func (c *FigmaClient) GetFileImages(fileKey string, nodeIDs []string) (map[string]string, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("file key cannot be empty")
	}
	if len(nodeIDs) == 0 {
		return nil, fmt.Errorf("node IDs cannot be empty")
	}

	// Clean the file key
	fileKey = c.extractFileKeyFromURL(fileKey)

	// Build the endpoint for getting images
	endpoint := fmt.Sprintf("%s/images/%s", c.baseURL, fileKey)

	// Add parameters
	params := url.Values{}
	params.Add("ids", strings.Join(nodeIDs, ","))
	params.Add("format", "png") // Default to PNG, could be configurable
	params.Add("scale", "1")    // Default scale, could be configurable
	endpoint += "?" + params.Encode()

	response, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file images from Figma API: %w", err)
	}

	var imageResponse struct {
		Images map[string]string `json:"images"`
		Err    string            `json:"err,omitempty"`
	}

	if err := json.Unmarshal(response, &imageResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Figma images response: %w", err)
	}

	if imageResponse.Err != "" {
		return nil, fmt.Errorf("Figma API error: %s", imageResponse.Err)
	}

	return imageResponse.Images, nil
}

// makeRequest is a helper method to make HTTP requests to Figma API
func (c *FigmaClient) makeRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers
	req.Header.Set("X-Figma-Token", c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "figma-parser-app/1.0")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.HandleAPIError(resp.StatusCode, responseBody)
	}

	return responseBody, nil
}

func (c *FigmaClient) extractFileKeyFromURL(input string) string {
	// If it's already just a file key (no slashes), return as-is
	if !strings.Contains(input, "/") {
		return input
	}

	// Handle full Figma URLs like:
	// https://www.figma.com/file/abc123/file-name
	// https://www.figma.com/design/abc123/file-name
	if strings.Contains(input, "figma.com") {
		parts := strings.Split(input, "/")
		for i, part := range parts {
			if (part == "file" || part == "design") && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}

	// If we can't extract it, return the original input
	return input
}

// ValidateAPIToken checks if the API token is valid by making a simple request
func (c *FigmaClient) ValidateAPIToken() error {
	endpoint := fmt.Sprintf("%s/me", c.baseURL)

	_, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("invalid API token: %w", err)
	}

	return nil
}

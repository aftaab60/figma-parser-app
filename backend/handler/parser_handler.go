package handler

import (
	"net/http"
	"parser-service/internal/errors"
	"parser-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ParserHandler struct {
	ParserService services.ParserService
}

func NewParserHandler(parserService services.ParserService) *ParserHandler {
	return &ParserHandler{
		ParserService: parserService,
	}
}

func (h *ParserHandler) ParseAndSaveFigmaFile(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		FigmaURL string `json:"figma_file_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Err:     "Invalid request body",
			Status:  http.StatusBadRequest,
			Message: "Please provide a valid Figma file URL",
		})
		return
	}

	savedFile, err := h.ParserService.ParseAndSaveFigmaFile(ctx, request.FigmaURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Err:     "Failed to parse Figma file",
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": savedFile})
}

func (h *ParserHandler) GetFigmaFileDetails(c *gin.Context) {
	ctx := c.Request.Context()

	// Get file ID from URL parameter
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Err:     "Invalid file ID",
			Status:  http.StatusBadRequest,
			Message: "File ID must be a valid number",
		})
		return
	}

	fileDetails, err := h.ParserService.GetFigmaFileWithDetails(ctx, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Err:     "Failed to get file details",
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fileDetails})
}

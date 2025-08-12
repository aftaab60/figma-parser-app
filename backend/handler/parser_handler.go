package handler

import (
	"net/http"
	"parser-service/internal/errors"
	"parser-service/services"

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

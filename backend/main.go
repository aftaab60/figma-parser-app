package main

import (
	"log"
	"parser-service/handler"
	"parser-service/internal/db_manager"
	"parser-service/internal/figma_manager"
	"parser-service/middlewares"
	"parser-service/repositories"
	"parser-service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // to recover from panics in execution

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:8080"}, // Allow both possible frontend ports
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	setupRoutes(r)
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	db := db_manager.InitPgsqlConnection()
	figmaManager := figma_manager.NewFigmaManager()
	figmaFilesRepo := repositories.NewFigmaFilesRepository(*db)
	componentsRepo := repositories.NewComponentsRepository(*db)
	instancesRepo := repositories.NewInstancesRepository(*db)
	parserHandler := handler.NewParserHandler(*services.NewParserService(figmaManager, figmaFilesRepo, componentsRepo, instancesRepo))

	// Apply middleware to routes that need Figma token validation
	r.POST("/parse-figma-file", middlewares.ValidateFigmaToken(figmaManager), parserHandler.ParseAndSaveFigmaFile)
	r.GET("/figma-files/:id", parserHandler.GetFigmaFileDetails) // No token needed for reading from DB
}

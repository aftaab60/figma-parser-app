package main

import (
	"log"
	"parser-service/handler"
	"parser-service/internal/db_manager"
	"parser-service/internal/figma_manager"
	"parser-service/repositories"
	"parser-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // to recover from panics in execution

	setupRoutes(r)
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	db := db_manager.InitPgsqlConnection()
	figmaManager := figma_manager.NewFigmaManager("<figma_api_token>") // to replace with actual
	figmaFilesRepo := repositories.NewFigmaFilesRepository(*db)
	componentsRepo := repositories.NewComponentsRepository(*db)
	instancesRepo := repositories.NewInstancesRepository(*db)

	parserHandler := handler.NewParserHandler(*services.NewParserService(figmaManager, figmaFilesRepo, componentsRepo, instancesRepo))
	r.POST("/parse-figma-file", parserHandler.ParseAndSaveFigmaFile)
}

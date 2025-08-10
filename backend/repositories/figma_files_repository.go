package repositories

import (
	"context"
	"parser-service/internal/db_manager"
	"parser-service/models"
)

// Interface to support dependency injection for Figma files repository
// just in case we want to switch to a different storage solution in the future.

type IFigmaFilesRepository interface {
	GetFigmaFileByID(ctx context.Context, id int64) (*models.FigmaFile, error)
	CreateFigmaFile(ctx context.Context, file *models.FigmaFile) (*models.FigmaFile, error)
	// Update(file *models.FigmaFile) (*models.FigmaFile, error) -- not in scope of this problem
	// Delete(id int64) error -- not in scope of this problem
}

type FigmaFilesRepository struct {
	DB db_manager.DB
}

func NewFigmaFilesRepository(db db_manager.DB) *FigmaFilesRepository {
	return &FigmaFilesRepository{DB: db}
}

func (r *FigmaFilesRepository) GetFigmaFileByID(ctx context.Context, id int64) (*models.FigmaFile, error) {
	query := "SELECT * FROM figma_files WHERE id = $1 AND active = TRUE"
	var file models.FigmaFile
	// scanning to individual columns for backward compatibility in future. This is better than using * and scanning to struct directly
	err := r.DB.GetRecord(ctx, query, id).Scan(&file.ID, &file.Name, &file.URL, &file.FileKey, &file.ImageURL, &file.Thumbnails, &file.CanvasWidth, &file.CanvasHeight, &file.ParsedAt, &file.CreatedAt, &file.UpdatedAt, &file.Active)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FigmaFilesRepository) CreateFigmaFile(ctx context.Context, file *models.FigmaFile) (*models.FigmaFile, error) {
	query := "INSERT INTO figma_files (name, url, file_key, image_url, thumbnails, canvas_width, canvas_height, parsed_at, created_at, updated_at, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), TRUE) RETURNING id, created_at, updated_at"
	err := r.DB.CreateRecord(ctx, query,
		file.Name,
		file.URL,
		file.FileKey,
		file.ImageURL,
		file.Thumbnails,
		file.CanvasWidth,
		file.CanvasHeight,
		file.ParsedAt).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return file, nil
}

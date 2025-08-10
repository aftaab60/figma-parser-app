package repositories

import (
	"context"
	"parser-service/internal/db_manager"
	"parser-service/models"
)

// Interface to support dependency injection for Components repository
// just in case we want to switch to a different storage solution in the future.

type IComponentsRepository interface {
	GetComponentByID(ctx context.Context, id int64) (*models.Component, error)
	GetComponentsByFigmaFileID(ctx context.Context, figmaFileID int64) ([]models.Component, error)
	CreateComponent(ctx context.Context, component *models.Component) (*models.Component, error)
}

type ComponentsRepository struct {
	DB db_manager.DB
}

func NewComponentsRepository(db db_manager.DB) *ComponentsRepository {
	return &ComponentsRepository{DB: db}
}

func (r *ComponentsRepository) GetComponentByID(ctx context.Context, id int64) (*models.Component, error) {
	query := "SELECT id, figma_file_id, node_id, name, type, description, x, y, width, height, z_index, properties, created_at, updated_at, active FROM components WHERE id = $1 AND active = TRUE"
	var component models.Component
	// scanning to individual columns for backward compatibility in future. This is better than using * and scanning to struct directly
	err := r.DB.GetRecord(ctx, query, id).Scan(
		&component.ID,
		&component.FigmaFileID,
		&component.NodeID,
		&component.Name,
		&component.Type,
		&component.Description,
		&component.X,
		&component.Y,
		&component.Width,
		&component.Height,
		&component.ZIndex,
		&component.Properties,
		&component.CreatedAt,
		&component.UpdatedAt,
		&component.Active)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *ComponentsRepository) GetComponentsByFigmaFileID(ctx context.Context, figmaFileID int64) ([]models.Component, error) {
	query := "SELECT id, figma_file_id, node_id, name, type, description, x, y, width, height, z_index, properties, created_at, updated_at, active FROM components WHERE figma_file_id = $1 AND active = TRUE ORDER BY z_index ASC"
	rows, err := r.DB.GetRecords(ctx, query, figmaFileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []models.Component
	for rows.Next() {
		var component models.Component
		err := rows.Scan(
			&component.ID,
			&component.FigmaFileID,
			&component.NodeID,
			&component.Name,
			&component.Type,
			&component.Description,
			&component.X,
			&component.Y,
			&component.Width,
			&component.Height,
			&component.ZIndex,
			&component.Properties,
			&component.CreatedAt,
			&component.UpdatedAt,
			&component.Active)
		if err != nil {
			return nil, err
		}
		components = append(components, component)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return components, nil
}

func (r *ComponentsRepository) CreateComponent(ctx context.Context, component *models.Component) (*models.Component, error) {
	query := "INSERT INTO components (figma_file_id, node_id, name, type, description, x, y, width, height, z_index, properties, created_at, updated_at, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW(), TRUE) RETURNING id, created_at, updated_at"
	err := r.DB.CreateRecord(ctx, query,
		component.FigmaFileID,
		component.NodeID,
		component.Name,
		component.Type,
		component.Description,
		component.X,
		component.Y,
		component.Width,
		component.Height,
		component.ZIndex,
		component.Properties).Scan(&component.ID, &component.CreatedAt, &component.UpdatedAt)
	if err != nil {
		return nil, err
	}
	component.Active = true
	return component, nil
}

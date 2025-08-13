package repositories

import (
	"context"
	"parser-service/internal/db_manager"
	"parser-service/models"
)

// Interface to support dependency injection for Instances repository
// just in case we want to switch to a different storage solution in the future.

type IInstancesRepository interface {
	GetInstanceByID(ctx context.Context, id int64) (*models.Instance, error)
	GetInstancesByComponentID(ctx context.Context, componentID int64) ([]models.Instance, error)
	GetInstancesByFigmaFileID(ctx context.Context, figmaFileID int64) ([]models.Instance, error)
	CreateInstance(ctx context.Context, instance *models.Instance) (*models.Instance, error)
}

type InstancesRepository struct {
	DB db_manager.DB
}

func NewInstancesRepository(db db_manager.DB) *InstancesRepository {
	return &InstancesRepository{DB: db}
}

func (r *InstancesRepository) GetInstanceByID(ctx context.Context, id int64) (*models.Instance, error) {
	query := "SELECT id, component_id, node_id, name, x, y, width, height, properties, created_at, updated_at, active FROM instances WHERE id = $1 AND active = TRUE"
	var instance models.Instance
	// scanning to individual columns for backward compatibility in future. This is better than using * and scanning to struct directly
	err := r.DB.GetRecord(ctx, query, id).Scan(
		&instance.ID,
		&instance.ComponentID,
		&instance.NodeID,
		&instance.Name,
		&instance.X,
		&instance.Y,
		&instance.Width,
		&instance.Height,
		&instance.Properties,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.Active)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *InstancesRepository) GetInstancesByComponentID(ctx context.Context, componentID int64) ([]models.Instance, error) {
	query := "SELECT id, component_id, node_id, name, x, y, width, height, properties, created_at, updated_at, active FROM instances WHERE component_id = $1 AND active = TRUE ORDER BY created_at ASC"
	rows, err := r.DB.GetRecords(ctx, query, componentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.Instance
	for rows.Next() {
		var instance models.Instance
		err := rows.Scan(
			&instance.ID,
			&instance.ComponentID,
			&instance.NodeID,
			&instance.Name,
			&instance.X,
			&instance.Y,
			&instance.Width,
			&instance.Height,
			&instance.Properties,
			&instance.CreatedAt,
			&instance.UpdatedAt,
			&instance.Active)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *InstancesRepository) GetInstancesByFigmaFileID(ctx context.Context, figmaFileID int64) ([]models.Instance, error) {
	query := `SELECT i.id, i.component_id, i.node_id, i.name, i.x, i.y, i.width, i.height, i.properties, i.created_at, i.updated_at, i.active 
			  FROM instances i 
			  INNER JOIN components c ON i.component_id = c.id 
			  WHERE c.figma_file_id = $1 AND i.active = TRUE AND c.active = TRUE 
			  ORDER BY i.created_at ASC`
	rows, err := r.DB.GetRecords(ctx, query, figmaFileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.Instance
	for rows.Next() {
		var instance models.Instance
		err := rows.Scan(
			&instance.ID,
			&instance.ComponentID,
			&instance.NodeID,
			&instance.Name,
			&instance.X,
			&instance.Y,
			&instance.Width,
			&instance.Height,
			&instance.Properties,
			&instance.CreatedAt,
			&instance.UpdatedAt,
			&instance.Active)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

func (r *InstancesRepository) CreateInstance(ctx context.Context, instance *models.Instance) (*models.Instance, error) {
	query := "INSERT INTO instances (component_id, node_id, name, x, y, width, height, properties, created_at, updated_at, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), TRUE) RETURNING id, created_at, updated_at"
	err := r.DB.CreateRecord(ctx, query,
		instance.ComponentID,
		instance.NodeID,
		instance.Name,
		instance.X,
		instance.Y,
		instance.Width,
		instance.Height,
		instance.Properties).Scan(&instance.ID, &instance.CreatedAt, &instance.UpdatedAt)
	if err != nil {
		return nil, err
	}
	instance.Active = true
	return instance, nil
}

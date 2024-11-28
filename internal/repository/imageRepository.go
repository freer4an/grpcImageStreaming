package repository

import (
	"context"
	"fmt"
	"github.com/freer4an/image-storage/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository struct {
	conn *pgxpool.Pool
}

func NewImageRepository(conn *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{conn}
}

func (r *ImageRepository) SaveImage(ctx context.Context, image *models.Image) (string, error) {
	const query = `INSERT INTO 
    		images ("id", "format", "width", "height", "original_path", "thumbnail_path") 
			VALUES ($1, $2, $3, $4, $5, $6)`

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction%w", err)
	}
	defer tx.Rollback(ctx)

	_, err = r.conn.Exec(ctx, query,
		image.Id,
		image.Format,
		image.Width,
		image.Height,
		image.OriginalPath,
		image.ThumbnailPath)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	return image.Id, nil
}

func (r *ImageRepository) GetImage(ctx context.Context, id string) (*models.Image, error) {
	const query = `SELECT original_path, thumbnail_path, width, height, uploaded_at FROM images WHERE id = $1`
	image := &models.Image{}
	row := r.conn.QueryRow(ctx, query, id)
	if err := row.Scan(&image.OriginalPath, &image.ThumbnailPath, &image.Width, &image.Height, &image.UploadedAt); err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	image.Id = id

	return image, nil
}

func (r *ImageRepository) DeleteImage(ctx context.Context, id string) error {
	const query = `DELETE FROM images WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

func (r *ImageRepository) ListImages(ctx context.Context) ([]models.Image, error) {
	const query = `SELECT id, format, width, height, original_path, thumbnail_path, uploaded_at FROM images`
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to exec query images: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Image])
}

func (r *ImageRepository) CountImages(ctx context.Context) (int, error) {
	const query = `SELECT COUNT(*) FROM images`
	var count int
	err := r.conn.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count images: %w", err)
	}
	return count, nil
}

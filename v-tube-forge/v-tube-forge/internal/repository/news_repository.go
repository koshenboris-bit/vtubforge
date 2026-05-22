package repository

import (
	"context"

	"backend-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	db *pgxpool.Pool
}

func NewNewsRepository(db *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) Create(ctx context.Context, item domain.News) (domain.News, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO news (title, content)
		VALUES ($1, $2)
		RETURNING id, title, content, created_at
	`, item.Title, item.Content).Scan(&item.ID, &item.Title, &item.Content, &item.CreatedAt)
	return item, err
}

func (r *NewsRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM news WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *NewsRepository) List(ctx context.Context) ([]domain.News, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, content, created_at
		FROM news
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.News
	for rows.Next() {
		var n domain.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

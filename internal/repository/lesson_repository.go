package repository

import (
	"context"

	"backend-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(db *pgxpool.Pool) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, item domain.Lesson) (domain.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO lessons (title, description, lesson_type, video_link)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, lesson_type, video_link, created_at
	`, item.Title, item.Description, item.LessonType, item.VideoLink).
		Scan(&item.ID, &item.Title, &item.Description, &item.LessonType, &item.VideoLink, &item.CreatedAt)
	return item, err
}

func (r *LessonRepository) Update(ctx context.Context, id int64, item domain.Lesson) (domain.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET title = $1, description = $2, lesson_type = $3, video_link = $4
		WHERE id = $5
		RETURNING id, title, description, lesson_type, video_link, created_at
	`, item.Title, item.Description, item.LessonType, item.VideoLink, id).
		Scan(&item.ID, &item.Title, &item.Description, &item.LessonType, &item.VideoLink, &item.CreatedAt)
	return item, err
}

func (r *LessonRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *LessonRepository) List(
	ctx context.Context,
	userID int64,
) ([]domain.Lesson, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			l.id,
			l.title,
			l.description,
			l.lesson_type,
			l.video_link,
			l.created_at,
			CASE
				WHEN ul.user_id IS NOT NULL THEN true
				ELSE false
			END as passed
		FROM lessons l
		LEFT JOIN user_lessons ul
			ON ul.lesson_id = l.id
			AND ul.user_id = $1
		ORDER BY l.id DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Lesson

	for rows.Next() {
		var l domain.Lesson

		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.LessonType,
			&l.VideoLink,
			&l.CreatedAt,
			&l.Passed,
		); err != nil {
			return nil, err
		}

		out = append(out, l)
	}

	return out, rows.Err()
}

func (r *LessonRepository) GetByID(ctx context.Context, id int64) (domain.Lesson, error) {
	var l domain.Lesson
	err := r.db.QueryRow(ctx, `
		SELECT id, title, description, lesson_type, video_link, created_at
		FROM lessons
		WHERE id = $1
	`, id).Scan(&l.ID, &l.Title, &l.Description, &l.LessonType, &l.VideoLink, &l.CreatedAt)
	return l, err
}

func (r *LessonRepository) MarkPassed(ctx context.Context, userID, lessonID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_lessons (user_id, lesson_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, lesson_id) DO NOTHING
	`, userID, lessonID)
	return err
}

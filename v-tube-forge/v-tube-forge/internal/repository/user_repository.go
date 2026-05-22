package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"backend-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, login, passwordHash string, role domain.Role) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (login, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, login, role, created_at
	`, login, passwordHash, role).Scan(&u.ID, &u.Login, &u.Role, &u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, string, error) {
	var u domain.User
	var passwordHash string
	err := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, role, created_at
		FROM users
		WHERE login = $1
	`, login).Scan(&u.ID, &u.Login, &passwordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return domain.User{}, "", err
	}
	return u, passwordHash, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, login, role, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Login, &u.Role, &u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, login, role, created_at
		FROM users
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Login, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetWithLastLesson(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, login, role, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Login, &u.Role, &u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	lesson, err := r.lastLesson(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	}
	if err == nil {
		u.LastLesson = &lesson
	}
	return u, nil
}

func (r *UserRepository) lastLesson(ctx context.Context, userID int64) (domain.Lesson, error) {
	var lesson domain.Lesson
	err := r.db.QueryRow(ctx, `
		SELECT l.id, l.title, l.description, l.lesson_type, l.video_link, l.created_at
		FROM user_lessons ul
		JOIN lessons l ON l.id = ul.lesson_id
		WHERE ul.user_id = $1
		ORDER BY ul.created_at DESC, l.id DESC
		LIMIT 1
	`, userID).Scan(&lesson.ID, &lesson.Title, &lesson.Description, &lesson.LessonType, &lesson.VideoLink, &lesson.CreatedAt)
	return lesson, err
}

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *RefreshTokenRepository) Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, hashToken(token), expiresAt)
	return err
}

func (r *RefreshTokenRepository) FindActiveByToken(ctx context.Context, token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow(ctx, `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
		ORDER BY id DESC
		LIMIT 1
	`, hashToken(token)).Scan(&userID)
	return userID, err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`, hashToken(token))
	return err
}

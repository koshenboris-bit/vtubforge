package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	Users   *UserRepository
	News    *NewsRepository
	Lessons *LessonRepository
	Tokens  *RefreshTokenRepository
}

func NewPostgresRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Users:   NewUserRepository(pool),
		News:    NewNewsRepository(pool),
		Lessons: NewLessonRepository(pool),
		Tokens:  NewRefreshTokenRepository(pool),
	}
}

package domain

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type LessonType string

const (
	LessonIntro LessonType = "intro"
	LessonFull  LessonType = "full"
)

type User struct {
	ID         int64     `json:"id"`
	Login      string    `json:"login"`
	Password   string    `json:"-"`
	Role       Role      `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	LastLesson *Lesson   `json:"lastLesson,omitempty"`
}

type Lesson struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	LessonType  LessonType `json:"lessonType"`
	VideoLink   string     `json:"videoLink"`
	CreatedAt   time.Time  `json:"createdAt"`
	Passed      bool       `json:"passed"`
}

type News struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserLesson struct {
	UserID    int64     `json:"userId"`
	LessonID  int64     `json:"lessonId"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type NewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type LessonRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	LessonType  LessonType `json:"lessonType"`
	VideoLink   string     `json:"videoLink"`
}

type PassLessonRequest struct {
	UserID int64 `json:"userId"`
}

type MeResponse struct {
	User User `json:"user"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

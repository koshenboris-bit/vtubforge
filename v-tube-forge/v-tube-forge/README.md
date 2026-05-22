# Backend Service

Go + PostgreSQL backend with:
- JWT access/refresh auth
- role-based access control
- handler/service/repository layers
- news, lessons, users, user-lessons

## Run with Docker

```bash
docker compose up --build
```

The app will start on `http://localhost:8080` and PostgreSQL on `localhost:5432`.

## Run locally

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/app?sslmode=disable
export JWT_ACCESS_SECRET=super-secret-access
export JWT_REFRESH_SECRET=super-secret-refresh
export HTTP_ADDR=:8080

go mod tidy
go run ./cmd/api
```

## Migrations

Migrations are applied automatically in the Docker entrypoint.

To apply them manually, run `migrations/001_init.sql` against PostgreSQL.

## Endpoints

Public:
- `POST /api/register`
- `POST /api/login`
- `POST /api/refresh`

Authenticated:
- `GET /api/news`
- `GET /api/lessons`
- `POST /api/lessons/{lessonId}/pass/{userId}`
- `GET /api/users/me`

Admin:
- `GET /api/users`
- `GET /api/users/{id}`
- `POST /api/news`
- `DELETE /api/news/{id}`
- `POST /api/lessons`
- `PUT /api/lessons/{id}`
- `DELETE /api/lessons/{id}`

## Notes

- Passwords are stored as bcrypt hashes and are never returned by the API.
- `/api/users/me` returns the user and the most recently passed lesson.
- The `POST /api/lessons/{lessonId}/pass/{userId}` endpoint allows the authenticated user to mark their own lesson as passed, while admins can mark any user.

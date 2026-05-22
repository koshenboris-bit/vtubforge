package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-service/internal/config"
	"backend-service/internal/handler"
	"backend-service/internal/middleware"
	"backend-service/internal/repository"
	"backend-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "backend-service/docs"
)

// @title Backend Service API
// @version 1.0
// @description Go + PostgreSQL backend service
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create db pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repos := repository.NewPostgresRepositories(pool)
	svcs := service.NewServices(cfg, repos)
	h := handler.NewHandler(cfg, svcs)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(api chi.Router) {
		api.Post("/register", h.Auth.Register)
		api.Post("/login", h.Auth.Login)
		api.Post("/refresh", h.Auth.Refresh)

		api.Group(func(pr chi.Router) {
			pr.Use(middleware.Auth(cfg.JWTAccessSecret, repos.Users))
			pr.Get("/users/me", h.Users.Me)
			pr.Get("/news", h.News.List)
			pr.Get("/lessons", h.Lessons.List)
			pr.Post("/lessons/{lessonId}/pass/{userId}", h.Lessons.Pass)
		})

		api.Group(func(ad chi.Router) {
			ad.Use(middleware.Auth(cfg.JWTAccessSecret, repos.Users))
			ad.Use(middleware.RequireRole("admin"))
			ad.Get("/users", h.Users.List)
			ad.Get("/users/{id}", h.Users.GetByID)

			ad.Post("/news", h.News.Create)
			ad.Delete("/news/{id}", h.News.Delete)

			ad.Post("/lessons", h.Lessons.Create)
			ad.Put("/lessons/{id}", h.Lessons.Update)
			ad.Delete("/lessons/{id}", h.Lessons.Delete)
		})
	})

	fs := http.FileServer(http.Dir("./web"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	r.Handle("/*", http.StripPrefix("/", fs))

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = srv.Shutdown(ctxShutdown)
	log.Println("server stopped")
}

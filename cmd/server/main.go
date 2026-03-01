package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GlebMoskalev/go-path-backend/content"
	"github.com/GlebMoskalev/go-path-backend/internal/config"
	"github.com/GlebMoskalev/go-path-backend/internal/database"
	"github.com/GlebMoskalev/go-path-backend/internal/handler"
	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/repository"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger, err := newLogger(cfg.Env)
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}
	defer logger.Sync()

	pool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		logger.Fatal("postgres pool error", zap.Error(err))
	}
	defer pool.Close()

	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal("redis client error", zap.Error(err))
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(pool)
	stateRepo := repository.NewStateRepository(redisClient)
	submissionRepo := repository.NewSubmissionRepository(pool)

	authService := service.NewAuthService(
		logger,
		userRepo,
		stateRepo,
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
		cfg.JWT.Secret,
		cfg.Google.UserInfoURL,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	userService := service.NewUserService(logger, userRepo)
	theoryService, err := service.NewTheoryService(content.TheoryFS, "theory", logger)
	if err != nil {
		logger.Fatal("failed to load theory", zap.Error(err))
	}

	taskService, err := service.NewTaskService(content.TasksFS, "tasks", logger)
	if err != nil {
		logger.Fatal("failed to load tasks", zap.Error(err))
	}

	sandboxService, err := service.NewSandboxService(logger, cfg.Sandbox.Image, cfg.Sandbox.Timeout, cfg.Sandbox.Memory)
	if err != nil {
		logger.Fatal("failed to create sandbox service", zap.Error(err))
	}
	submissionService := service.NewSubmissionService(logger, taskService, sandboxService, submissionRepo)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	theoryHandler := handler.NewTheoryHandler(theoryService)
	taskHandler := handler.NewTaskHandler(taskService, submissionService)

	authMiddleware := middleware.NewAuthMiddleware(authService, userService)

	router := chi.NewRouter()

	router.Route("/api", func(api chi.Router) {
		api.Get("/google/login", authHandler.GoogleLogin)
		api.Get("/google/callback", authHandler.GoogleCallback)
		api.Post("/refresh", authHandler.RefreshToken)

		api.Route("/auth", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Post("/logout", authHandler.Logout)
		})

		api.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/profile", userHandler.GetProfile)
			r.Put("/profile", userHandler.UpdateProfile)
			r.Delete("/account", userHandler.DeleteAccount)
		})

		api.Route("/theory", func(r chi.Router) {
			r.Get("/", theoryHandler.ListChapter)
			r.Get("/{chapterSlug}", theoryHandler.GetChapter)
			r.Get("/{chapterSlug}/{lessonSlug}", theoryHandler.GetLesson)
		})

		api.Route("/tasks", func(r chi.Router) {
			r.Get("/", taskHandler.ListChapters)
			r.Get("/{chapterSlug}", taskHandler.GetChapter)
			r.Get("/{chapterSlug}/{taskSlug}", taskHandler.GetTask)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Post("/{chapterSlug}/{taskSlug}/submit", taskHandler.Submit)
			})
		})
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	server := &http.Server{Addr: addr, Handler: router}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	sig := <-quit
	logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}

func newLogger(env string) (*zap.Logger, error) {
	if env == "prod" {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

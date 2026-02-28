package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

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
		if err := http.ListenAndServe(addr, router); err != nil {
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

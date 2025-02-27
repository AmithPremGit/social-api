package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"social-api/internal/store"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"

	"social-api/internal/auth"
	"social-api/internal/cache"
	"social-api/internal/config"
	"social-api/internal/db"
	"social-api/internal/handler"
	appMiddleware "social-api/internal/middleware"

	"social-api/internal/store/postgres"
)

// Application version
const version = "1.0.0"

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Printf("Starting social-api version %s", version)

	// Load configuration
	cfg := config.Load()
	logger.Printf("Environment: %s", cfg.Env)

	// Connect to database
	database, err := db.New(cfg.DB)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()
	logger.Printf("Connected to database on %s", cfg.DB.Host)

	// Initialize Redis client if enabled
	var redisClient *redis.Client
	var cacheService cache.Cache
	if cfg.Redis.Enabled {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		// Test Redis connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			logger.Fatalf("Redis connection failed: %v", err)
		}

		cacheService = cache.NewRedisCache(redisClient)
		logger.Printf("Connected to Redis on %s", cfg.Redis.Addr)

		defer redisClient.Close()
	} else {
		logger.Println("Redis cache disabled")
	}

	// Initialize stores
	userStore := postgres.NewUserStore(database)
	postStore := postgres.NewPostStore(database)

	// Initialize authenticator
	authenticator := auth.NewJWTAuthenticator(
		cfg.Auth.TokenSecret,
		cfg.Auth.TokenIssuer,
		cfg.Auth.TokenAudience,
	)

	// Initialize rate limiter
	rateLimiter := appMiddleware.NewFixedWindowRateLimiter(
		cfg.RateLimiter.RequestsPerWindow,
		cfg.RateLimiter.WindowDuration,
	)

	// Initialize application handler
	app := handler.NewApplication(
		cfg,
		logger,
		authenticator,
		cacheService,
		userStore,
		postStore,
	)

	// Set up router with middleware
	router := setupRouter(app, cfg, logger, authenticator, userStore, rateLimiter)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("Starting server on port %s", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received or an error occurs
	select {
	case err := <-serverErrors:
		logger.Fatalf("Server error: %v", err)
	case <-shutdown:
		logger.Println("Starting graceful shutdown")

		// Create a context with a timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Printf("Graceful shutdown error: %v", err)
			err = srv.Close()
			if err != nil {
				logger.Fatalf("Forced shutdown error: %v", err)
			}
		}

		// Check for errors during shutdown
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Println("Shutdown timeout: some connections may have been dropped")
		}

		logger.Println("Server shutdown complete")
	}
}

// setupRouter configures the router with middleware and routes
func setupRouter(
	app *handler.Application,
	cfg config.Config,
	logger *log.Logger,
	authenticator *auth.JWTAuthenticator,
	userStore store.UserStore,
	rateLimiter appMiddleware.RateLimiter,
) http.Handler {
	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Use our custom logging middleware instead of Chi's built-in logger
	r.Use(appMiddleware.LoggingMiddleware(logger))

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Apply rate limiting if enabled
	if cfg.RateLimiter.Enabled {
		r.Use(appMiddleware.RateLimiterMiddleware(rateLimiter, logger))
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", app.HealthCheck)
		r.Post("/users", app.RegisterUser)
		r.Post("/auth/token", app.CreateToken)

		// Protected routes
		r.Group(func(r chi.Router) {
			// Apply auth middleware
			r.Use(auth.Middleware(authenticator, userStore))

			// User routes
			r.Get("/users/me", app.GetCurrentUser)

			// Post routes
			r.Route("/posts", func(r chi.Router) {
				r.Get("/", app.ListPosts)
				r.Post("/", app.CreatePost)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetPost)
					r.Put("/", app.UpdatePost)
					r.Delete("/", app.DeletePost)
				})
			})
		})
	})

	return r
}

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/provider"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/middleware"
)

type HttpServer struct {
	cfg    *config.Config
	server *http.Server
}

func NewHttpServer(cfg *config.Config) *HttpServer {
	gin.SetMode(cfg.App.Environment)

	router := gin.New()
	router.ContextWithFallback = true
	router.HandleMethodNotAllowed = true

	RegisterMiddleware(router, cfg)
	provider.BootstrapHttp(cfg, router)

	return &HttpServer{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
			Handler: router,
		},
	}
}

func (s *HttpServer) Start() {
	log.Logger.Info("Running HTTP server on port:", s.cfg.HttpServer.Port)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Fatal("Error while HTTP server listening:", err)
	}
	log.Logger.Info("HTTP server is not receiving new requests...")
}

func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.HttpServer.GracePeriod)*time.Second)
	defer cancel()

	log.Logger.Info("Attempting to shut down the HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Logger.Fatal("Error shutting down HTTP server:", err)
	}
	log.Logger.Info("HTTP server shut down gracefully")
}

func RegisterMiddleware(router *gin.Engine, cfg *config.Config) {
	middlewares := []gin.HandlerFunc{
		gin.Recovery(),
		gzip.Gzip(gzip.BestSpeed),
		middleware.Logger(log.Logger),
		middleware.ErrorHandler(),
		middleware.RequestTimeout(cfg.HttpServer.RequestTimeoutPeriod),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowAllOrigins:  true,
			AllowCredentials: true,
		}),
	}

	router.Use(middlewares...)
}

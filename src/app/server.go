package app

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
	db     *mongo.Database
}

type ServerConfig interface {
	GetHTTPAddr() string
}

// NewAPI create a new WebHookAPI with Connection, logger and router
func NewServer(cfg ServerConfig, logger *zap.Logger, routerBuilder *PipelineBuilder, db *mongo.Database) *Server {
	server := &Server{
		server: &http.Server{
			Addr:    cfg.GetHTTPAddr(),
			Handler: routerBuilder.AddMiddleware().AddRoutes().Build(),
		},
		logger: logger,
		db:     db,
	}

	return server
}

// Listen app and block on this method until os.Interrupt or os.Kill received
func (s *Server) Listen(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Couldn't listen: %+v\n", zap.Error(err))
		}
	}()
	<-ctx.Done()
	s.logger.Info("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	s.server.SetKeepAlivesEnabled(false)
	// ignore error since it will be "Err shutting down server : context canceled"
	return s.server.Shutdown(shutdownCtx)
}

type PipelineBuilder struct {
	router *gin.Engine
	logger *zap.Logger
	routes []RouterBinder
}

func NewPipelineBuilder(logger *zap.Logger, routes []RouterBinder) *PipelineBuilder {
	return &PipelineBuilder{router: gin.New(), logger: logger, routes: routes}
}

type RouterBinder interface {
	// BindRoutes adds some routes to gin engine
	BindRoutes(*gin.Engine)
}

func (b *PipelineBuilder) AddMiddleware() *PipelineBuilder {
	b.router.Use(ginzap.Ginzap(b.logger, time.RFC3339, true))
	b.router.Use(ginzap.RecoveryWithZap(b.logger, true))
	b.router.Use(defaultHeaders())
	b.router.Use(cors.Default())
	return b
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.Next()
	}
}

func (b *PipelineBuilder) AddRoutes() *PipelineBuilder {
	for _, r := range b.routes {
		r.BindRoutes(b.router)
	}
	return b
}

func (b *PipelineBuilder) Build() *gin.Engine {
	return b.router
}

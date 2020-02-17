package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

type ServerConfig interface {
	GetHTTPAddr() string
}

// NewAPI create a new WebHookAPI with DB, logger and router
func NewServer(cfg ServerConfig, logger *zap.Logger, routerBuilder *PipelineBuilder) *Server {
	server := &Server{
		server: &http.Server{
			Addr:    cfg.GetHTTPAddr(),
			Handler: routerBuilder.AddMiddleware().AddRoutes().Build(),
		},
		logger: logger,
	}

	return server
}

// Run server and block on this method until os.Interrupt or os.Kill received
func (s *Server) Run() {
	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		s.logger.Info("Shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Fatal("Couldn't gracefully shutdown the server.go: %+v\n", zap.Error(err))
		}
		close(done)
	}()

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Couldn't listen: %+v\n", zap.Error(err))
	}

	<-done
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

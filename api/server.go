package api

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/config"
)

const defaultServerShutdownTimeout = 5 * time.Second

// Server implements http server.
type Server struct {
	server *http.Server
}

// NewServer creates (but doesn't start) new instance of http server.
func NewServer(cfg config.Config, routes ...RouteBinder) *Server {
	// defaulting to gin release mode because we do not need additional formatting.
	gin.SetMode(gin.ReleaseMode)

	server := &Server{
		server: &http.Server{
			Addr:    ":" + cfg.WebAPIPort,
			Handler: newPipelineBuilder(routes, cfg).addMiddleware().addRoutes().build(),
		},
	}
	return server
}

// Start web server.
func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start API server")
		}
	}()
}

// Shutdown web server and kill all keep-alive connections.
func (s *Server) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, defaultServerShutdownTimeout)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	// we do not care about this particular error because we can't do anything when we already shutting down.
	_ = s.server.Shutdown(shutdownCtx)
}

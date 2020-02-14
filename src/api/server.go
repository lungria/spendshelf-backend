package api

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	router *gin.Engine
}

type RouterBinder interface {
	// BindRoutes adds some routes to gin engine
	BindRoutes(*gin.Engine)
}

type ServerConfig interface {
	GetHTTPAddr() string
}

type Binder struct {
	All []RouterBinder
}

// NewAPI create a new WebHookAPI with DB, logger and router
func NewAPI(cfg ServerConfig, logger *zap.Logger, bt []byte) (*Server, error) {
	fmt.Printf("%s", bt)
	server := &Server{
		&http.Server{
			Addr: cfg.GetHTTPAddr(),
		},
		configureMiddleware(logger),
	}
	server.bindRoutes(&Binder{})
	return server, nil
}

func configureMiddleware(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(defaultHeaders())
	router.Use(cors.Default())
	return router
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.Next()
	}
}

func (s *Server) bindRoutes(binders *Binder) {
	for _, v := range binders.All {
		v.BindRoutes(s.router)
	}
}

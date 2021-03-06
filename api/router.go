package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/app/config"
	cors "github.com/rs/cors/wrapper/gin"
)

// RouteBinder abstracts some entity (usually handler) that can register it's routes in gin router.
type RouteBinder interface {
	// BindRoutes adds some routes to gin engine
	BindRoutes(*gin.Engine)
}

type pipelineBuilder struct {
	router *gin.Engine
	routes []RouteBinder
	cfg    config.Config
}

func newPipelineBuilder(routes []RouteBinder, cfg config.Config) *pipelineBuilder {
	return &pipelineBuilder{router: gin.New(), routes: routes, cfg: cfg}
}

// addMiddleware builds middleware pipeline.
func (b *pipelineBuilder) addMiddleware() *pipelineBuilder {
	b.router.Use(gin.RecoveryWithWriter(log.Writer())).
		Use(defaultHeaders()).
		Use(cors.New(cors.Options{
			AllowedOrigins: []string{b.cfg.CORSHost},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			MaxAge: int((24 * time.Hour).Seconds()),
		}))

	return b
}

// addRoutes binds all routes.
func (b *pipelineBuilder) addRoutes() *pipelineBuilder {
	for _, r := range b.routes {
		r.BindRoutes(b.router)
	}

	return b
}

// Builds gin middleware pipeline.
func (b *pipelineBuilder) build() *gin.Engine {
	return b.router
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

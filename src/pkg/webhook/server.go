package webhook

import (
	"github.com/lungria/spendshelf-backend/src/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	HTTPServer *http.Server
	MongoDB    *mongo.Database
	Logger 		zap.SugaredLogger
}

// NewServer is create a new Server
func NewServer(addr, dbname, MongoURI string) *Server {
	s := &Server{
		HTTPServer: nil,
		MongoDB:    nil,
	}
	s.newDatabase(dbname, MongoURI)
	s.newHTTPServer(addr)

	return s
}

// NewServer is create a new HTTP server
func (s *Server) newHTTPServer(addr string) {
	router := gin.New()

	logger, _ := zap.NewProduction()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", s.WebHookHandler)

	s.HTTPServer = &http.Server{
		Addr:              addr,
		Handler:           router,
	}
	s.Logger = *logger.Sugar()
}

// NewDatabase is create a new database connection
func (s *Server) newDatabase(dbname, URI string) {
	conn, err := db.Connect(dbname, URI)
	if err != nil {
		s.Logger.Fatalw("Connection to mongo failed", "Database", dbname, "Address", URI, "Error", err)
		return
	}
	s.MongoDB = conn
}

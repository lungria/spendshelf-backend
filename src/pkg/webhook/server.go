package webhook

import (
	"github.com/lungria/spendshelf-backend/src/pkg/webhook/api"
	"go.uber.org/zap"
)

type Server struct {
	WebHookAPI  *api.WebHookAPI
	Logger 		*zap.SugaredLogger
}

// NewServer is create a new Server
func NewServer(addr, dbname, mongoURI string) (*Server, error) {
	webHookAPI, err := api.NewAPI(addr, dbname, mongoURI)
	if err != nil {
		return nil, err
	}
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &Server{
		WebHookAPI: webHookAPI,
		Logger:     logger.Sugar(),
	}, nil
}



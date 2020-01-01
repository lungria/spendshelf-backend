package config

import (
	"net/http"

	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

package config

import (
	"net/http"

	"go.uber.org/zap"
)

// Dependencies is struct which defined dependencies for Wire
type Dependencies struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

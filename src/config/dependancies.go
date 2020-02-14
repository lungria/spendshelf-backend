package config

import (
	"github.com/lungria/spendshelf-backend/src/api"
	"go.uber.org/zap"
)

// Dependencies is struct which defined dependencies for Wire
type Dependencies struct {
	Logger *zap.SugaredLogger
	Server *api.Server
}

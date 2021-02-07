package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/app/config"
	"github.com/lungria/spendshelf-backend/app/job"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/rs/zerolog/log"
)

// State stores information about app dependencies and allows to manage it's lifecycle.
type State struct {
	API       *api.Server
	Scheduler *job.Scheduler
	DB        *pgxpool.Pool
	Importer  *importer.Importer
	Config    config.Config
}

// Close releases all resources and stops all background jobs.
func (s *State) Close() {
	s.API.Shutdown()
	s.Scheduler.Wait()
	s.DB.Close()
	log.Info().Msg("app shutdown finished gracefully")
}

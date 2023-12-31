package app

import (
	"fmt"
	"net/http"

	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/services/waitlist"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	Config          *config.Config
	WaitlistService *waitlist.Service
	Routes          *chi.Mux
}

func readConfigs() *config.Config {
	config, err := config.Read("./config.yml")
	if err != nil {
		log.Fatal().Err(err)
	}

	return config
}

func createServices(conf *config.Config) *waitlist.Service {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	waitlistService := waitlist.New(conf.Database)

	return waitlistService
}

func new() *application {
	config := readConfigs()

	waitlistService := createServices(config)
	routes := createRoutes(waitlistService)

	app := application{
		Config:          config,
		WaitlistService: waitlistService,
		Routes:          routes,
	}

	return &app
}

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := new()

	addr := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)
	log.Info().Msg(fmt.Sprintf("Running application at %s", addr))
	http.ListenAndServe(addr, app.Routes)
}

package http

import (
	"github.com/khivuksergey/portmonetka.wallet/config"
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/gorm"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/service"
	"github.com/khivuksergey/webserver"
	"github.com/khivuksergey/webserver/logger"
)

const configPath = "config.json"

func NewServer() webserver.Server {
	config.LoadEnv()

	cfg := config.LoadConfiguration(configPath)

	db := gorm.NewDbManager(cfg.DB)

	services := service.NewServiceManager(db.InitRepositoryManager())

	log := logger.Default.SetLevel(logger.GetLogLevelFromString(cfg.Logger.LogLevel))

	router := NewRouter(cfg, services, log)

	server := webserver.
		NewServer(router).
		WithConfig(&cfg.Server).
		AddLogger(log).
		AddStopHandlers(webserver.NewStopHandler("Database", db.Close))

	return server
}

package main

import (
	_ "github.com/mattn/go-sqlite3"
	"grafana-matrix-alerts/config"
	"grafana-matrix-alerts/logger"
	"grafana-matrix-alerts/mautrixClient"
	"grafana-matrix-alerts/router"
)

func main() {
	// Init logger
	logger.InitLogger(config.DebugLog)

	// Init gin router
	go func() {
		_ = router.WebhookEndpoints().Run(":" + config.Port)
	}()

	// Init matrix connection
	mautrixClient.Connect()

	logger.Log.Info().Msg("App started")

	select {}
}

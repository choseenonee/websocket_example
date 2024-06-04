package main

import (
	"websockets/internal/delivery"
	"websockets/internal/metrics"
	"websockets/pkg/config"
	"websockets/pkg/database"
	"websockets/pkg/log"
)

func main() {
	config.InitConfig()

	logger, errFile, infoFile := log.InitLogger()

	defer errFile.Close()

	defer infoFile.Close()

	db := database.MustGetDB()

	sendMessageMetrics := metrics.InitMetrics()

	delivery.Start(logger, db, sendMessageMetrics)
}

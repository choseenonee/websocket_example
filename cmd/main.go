package main

import (
	"fmt"
	"github.com/spf13/viper"
	"websockets/internal/delivery"
	"websockets/internal/metrics"
	"websockets/pkg/config"
	"websockets/pkg/database"
	"websockets/pkg/log"
	"websockets/pkg/tracing"
)

const serviceName = "wsapp"

func main() {
	logger, errFile, infoFile := log.InitLogger()
	logger.Info("Logger initialized!")

	defer errFile.Close()

	defer infoFile.Close()

	config.InitConfig()
	logger.Info("Config initialized!")

	jaegerURL := fmt.Sprintf("http://%v:%v/api/traces", viper.GetString(config.JaegerHost), viper.GetString(config.JaegerPort))
	tracer := tracing.InitTracer(jaegerURL, serviceName)
	logger.Info("Tracing initialized!")

	db := database.MustGetDB()
	logger.Info("Database initialized!")

	sendMessageMetrics := metrics.InitMetrics()
	logger.Info("Metrics initialized!")

	delivery.Start(logger, db, sendMessageMetrics, tracer)
}

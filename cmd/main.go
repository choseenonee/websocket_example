package main

import (
	"websockets/internal/delivery"
	"websockets/pkg/log"
)

func main() {
	logger, errFile, infoFile := log.InitLogger()

	defer errFile.Close()

	defer infoFile.Close()

	delivery.Start(logger)
}

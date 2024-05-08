package main

import (
	"websockets/internal"
	"websockets/pkg/log"
)

func main() {
	logger, errFile, infoFile := log.InitLogger()

	defer errFile.Close()

	defer infoFile.Close()

	internal.Start(logger)
}

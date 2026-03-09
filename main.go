package main

import (
	"log/slog"
	"os"

	"github.com/SealJonny/dyndns-go/server"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting dyndns-go", "version", version)

	zoneID, exists := os.LookupEnv("CF_ZONE_ID")
	if !exists {
		slog.Error("CF_ZONE_ID is not set")
		os.Exit(1)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	server := server.New(port, zoneID)
	slog.Info("starting server", "port", port)
	err := server.Start()
	if err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}

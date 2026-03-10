package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/SealJonny/dyndns-go/notification"
	"github.com/SealJonny/dyndns-go/server"
)

var version = "dev"

func checkENV(name, msg string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		slog.Error(msg)
		os.Exit(1)
	}

	return value
}

func setupNotification() {
	smtpEnabled := false
	stmpEnabledStr, exists := os.LookupEnv("SMTP_ENABLE")
	if exists {
		var err error
		smtpEnabled, err = strconv.ParseBool(stmpEnabledStr)
		if err != nil {
			slog.Error("SMTP_ENABLE is not a valid boolean")
			os.Exit(1)
		}
	}

	if !smtpEnabled {
		return
	}

	host := checkENV("SMTP_HOST", "SMTP is enabled but SMTP_HOST is not set")
	port := checkENV("SMTP_PORT", "SMTP is enabled but SMTP_PORT is not set")
	user := checkENV("SMTP_USERNAME", "SMTP is enabled but SMTP_USERNAME is not set")
	password := checkENV("SMTP_PASSWORD", "SMTP is enabled but SMTP_PASSWORD is not set")
	receiver := checkENV("SMTP_RECEIVER", "SMTP is enabled but SMTP_RECEIVER is not set")

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		slog.Error("SMTP_PORT is not an integer")
		os.Exit(1)
	}

	smtpClient := notification.New(host, parsedPort, user, password, receiver)
	notification.SetSMTPDefault(smtpClient)
	notification.EnableSMTP()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting dyndns-go", "version", version)

	setupNotification()

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

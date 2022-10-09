package main

import (
	"os"
	"strings"

	"github.com/litestack-hq/lgst-tracker/cmd"
	"github.com/litestack-hq/lgst-tracker/helpers/config"

	"github.com/rs/zerolog"
)

func main() {
	config := config.New()
	logger := setupLogger(config.APP_NAME, config.APP_ENV)

	cmd.Run(config, logger)
}

func setupLogger(serviceName string, env string) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp()
	logger = logger.Str("service-name", serviceName)

	if strings.ToLower(env) == "production" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		return logger.Logger()
	}

	if strings.ToLower(env) == "development" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		return logger.Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger.Logger()
}

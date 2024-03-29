package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/litestack-hq/lgst-common/helpers/utils"
	"github.com/litestack-hq/lgst-tracker/helpers/config"
	"github.com/litestack-hq/lgst-tracker/tcp/device_ping_module"
	"github.com/rs/zerolog"

	app_http "github.com/litestack-hq/lgst-tracker/http"
	"github.com/litestack-hq/lgst-tracker/tcp"
)

func Run(conf *config.Config, logger zerolog.Logger) {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		default:
			logger.Panic().Msg("command not found")
		}
		return
	}

	runApp(conf, logger)
}

func runApp(conf *config.Config, logger zerolog.Logger) {
	if conf.PORT == 0 {
		logger.Warn().Msg("HTTP port not configured: generating a random port")
		port, err := utils.GetFreePort()

		if err != nil {
			logger.Panic().Err(err).Msg("failed to generate a free port")
		}

		conf.PORT = port
	}

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", conf.PORT),
		Handler: app_http.Handler(app_http.HandlerOpts{
			Conf:   conf,
			Logger: logger,
		}),
	}

	closeChannel := make(chan struct{})

	devicePingModule := device_ping_module.New(tcp.ModuleOpts{
		Name:        "Device ping module",
		Protocol:    tcp.PROTOCOL_TCP,
		DefaultPort: "7000",
		Logger:      logger,
	})

	go devicePingModule.Run()

	go func() {
		sigInt := make(chan os.Signal, 1)
		signal.Notify(sigInt, os.Interrupt)
		<-sigInt

		logger.Info().Msg("shutting down HTTP server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Err(err).Msg("HTTP server shutdown error")
		}

		logger.Info().Msg("shutting down device ping module")
		if err := devicePingModule.Close(); err != nil {
			logger.Err(err).Msg("device ping module shutdown error")
		}

		close(closeChannel)
	}()

	logger.Info().Msgf("HTTP server running on port %d", conf.PORT)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Panic().Err(err).Msg("HTTP server listen and serve failed")
	}

	<-closeChannel
}

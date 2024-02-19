package main

import (
	"context"
	"flag"

	"github.com/jackvonhouse/enrichment/app"
	"github.com/jackvonhouse/enrichment/config"
	"github.com/jackvonhouse/enrichment/pkg/log"
	"github.com/jackvonhouse/enrichment/pkg/shutdown"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogrusLogger()

	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"config/config.toml",
		"The path to the configuration file",
	)

	flag.Parse()

	config, err := config.New(configPath, logger)
	if err != nil {
		logger.Error(err)

		return
	}

	app, err := app.New(ctx, config, logger)
	if err != nil {
		logger.Error(err)

		return
	}

	go app.Run()

	shutdown.Graceful(ctx, cancel, logger, app)
}

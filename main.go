package main

import (
	"context"
	"flag"
	"publish_to_messaging_service/internal/app"
	"publish_to_messaging_service/internal/config"
	"publish_to_messaging_service/pkg/logger"
)

var (
	mode    = flag.String("mode", "dev", "Deployment mode")
	runMode = flag.String("run", "test", "Run mode")
	ctx     = context.Background()
)

func parseArgs() (string, string) {
	flag.Parse()
	return *mode, *runMode
}

func init() {
	*mode, *runMode = parseArgs()
	ctx = context.Background()

	if err := config.InitConfig(*mode, *runMode); err != nil {
		panic(err)
	}

	if err := logger.InitLogger(&config.ModeConfig.Logging, *mode); err != nil {
		panic(err)
	}
}

func main() {
	defer logger.Sugar.Sync()

	app.Application(ctx, *mode, *runMode)
}

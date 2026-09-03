package main

import (
	"context"
	"flag"
	"os"

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/runner"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", ".config/gildedrose.yml", "Path to YAML configuration file")
	flag.Parse()

	cfg, err := config.LoadYAML(configFile)
	if err != nil {
		logger.LogError("config load failed: %v", err)
		os.Exit(1)
	}

	finalState, err := runner.Run(context.Background(), cfg)
	if err != nil {
		logger.LogError("run failed: %v", err)
		os.Exit(1)
	}
	if !runner.Success(finalState) {
		os.Exit(1)
	}
}

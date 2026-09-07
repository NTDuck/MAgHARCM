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
	flag.StringVar(&configFile, "config", "", "Path to YAML configuration file (required)")
	flag.Parse()

	if configFile == "" {
		// If omitted, check for default request file or prompt user
		configFile = ".config/gildedrose.yml"
	}

	cfg := config.MustLoadYAML(configFile)

	finalState, err := runner.Run(context.Background(), cfg)
	if err != nil {
		logger.LogError("run failed: %v", err)
		os.Exit(1)
	}
	if !runner.Success(finalState) {
		os.Exit(1)
	}
}

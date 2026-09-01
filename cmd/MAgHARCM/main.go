package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/runner"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yml", "Path to YAML configuration file")
	flag.Parse()

	cfg, err := config.LoadYAML(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load YAML configuration from `%s`: %v\n", configFile, err)
		os.Exit(1)
	}

	finalState, err := runner.Run(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !runner.Success(finalState) {
		os.Exit(1)
	}
}

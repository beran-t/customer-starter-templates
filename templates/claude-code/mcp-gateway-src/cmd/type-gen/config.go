package main

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds the configuration for the type generator
type Config struct {
	CatalogPaths  []string
	SpecOutput    string
	MappingOutput string
	Verbose       bool
}

// ParseFlags parses command-line flags and returns a Config
func ParseFlags() *Config {
	cfg := &Config{}

	var catalogPathsStr string
	flag.StringVar(&catalogPathsStr, "catalogs", "./docker-catalog.yaml", "Comma-separated list of catalog file paths")
	flag.StringVar(&cfg.SpecOutput, "spec-output", "spec.json", "Output path for generated JSON schema")
	flag.StringVar(&cfg.MappingOutput, "mapping-output", "mapping.json", "Output path for generated instruction mapping")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	// Split catalog paths
	if catalogPathsStr != "" {
		cfg.CatalogPaths = strings.Split(catalogPathsStr, ",")
		for i := range cfg.CatalogPaths {
			cfg.CatalogPaths[i] = strings.TrimSpace(cfg.CatalogPaths[i])
		}
	}

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.CatalogPaths) == 0 {
		return fmt.Errorf("at least one catalog path must be specified")
	}
	if c.SpecOutput == "" {
		return fmt.Errorf("spec output path cannot be empty")
	}
	if c.MappingOutput == "" {
		return fmt.Errorf("mapping output path cannot be empty")
	}
	return nil
}

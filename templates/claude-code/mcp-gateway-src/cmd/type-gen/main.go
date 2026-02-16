package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Parse configuration
	cfg := ParseFlags()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	var logger *zap.Logger
	var err error
	if cfg.Verbose {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	logger.Info("starting type generation",
		zap.Strings("catalogs", cfg.CatalogPaths),
		zap.String("spec_output", cfg.SpecOutput),
		zap.String("mapping_output", cfg.MappingOutput))

	// Read catalog
	cat, err := catalog.ReadFrom(ctx, cfg.CatalogPaths)
	if err != nil {
		logger.Fatal("failed to read catalog", zap.Error(err))
	}
	logger.Info("catalog loaded",
		zap.Int("servers", len(cat.Servers)))

	// Transform catalog
	transformer := NewCatalogTransformer(logger)
	configs, requireds, instructionMap, metadataMap := transformer.Transform(cat)

	// Generate combined schema
	schemaGen := NewSchemaGenerator(logger)
	combinedSchema, err := schemaGen.GenerateCombinedSchema(configs, requireds, metadataMap)
	if err != nil {
		logger.Fatal("failed to generate combined schema", zap.Error(err))
	}

	// Convert schema to pretty JSON
	prettySchema, err := schemaGen.ToPrettyJSON(combinedSchema)
	if err != nil {
		logger.Fatal("failed to prettify schema", zap.Error(err))
	}

	// Write schema to spec.json
	if err := os.WriteFile(cfg.SpecOutput, prettySchema, 0o644); err != nil {
		logger.Fatal("failed to write spec file",
			zap.String("path", cfg.SpecOutput),
			zap.Error(err))
	}
	logger.Info("spec file written",
		zap.String("path", cfg.SpecOutput),
		zap.Int("bytes", len(prettySchema)))

	// Write instruction map to mapping.json
	mappingJSON, err := json.MarshalIndent(instructionMap, "", "  ")
	if err != nil {
		logger.Fatal("failed to marshal instruction map", zap.Error(err))
	}

	if err := os.WriteFile(cfg.MappingOutput, mappingJSON, 0o644); err != nil {
		logger.Fatal("failed to write mapping file",
			zap.String("path", cfg.MappingOutput),
			zap.Error(err))
	}
	logger.Info("mapping file written",
		zap.String("path", cfg.MappingOutput),
		zap.Int("bytes", len(mappingJSON)))

	logger.Info("type generation complete")
}

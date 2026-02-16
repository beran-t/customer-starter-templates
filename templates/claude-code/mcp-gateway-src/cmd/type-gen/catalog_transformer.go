package main

import (
	"fmt"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"go.uber.org/zap"
)

// CatalogTransformer transforms a catalog into configuration maps
type CatalogTransformer struct {
	logger             *zap.Logger
	instructionBuilder *InstructionBuilder
}

// NewCatalogTransformer creates a new CatalogTransformer
func NewCatalogTransformer(logger *zap.Logger) *CatalogTransformer {
	return &CatalogTransformer{
		logger:             logger,
		instructionBuilder: NewInstructionBuilder(logger),
	}
}

// Transform transforms a Catalog to configuration maps
func (t *CatalogTransformer) Transform(cat catalog.Catalog) (ConfigMap, RequiredMap, InstructionMap, MetadataMap) {
	// Early return for empty catalog
	if len(cat.Servers) == 0 {
		t.logger.Warn("empty catalog provided")
		return ConfigMap{}, RequiredMap{}, InstructionMap{}, MetadataMap{}
	}

	configs := make(ConfigMap, len(cat.Servers))
	requireds := make(RequiredMap, len(cat.Servers))
	instructionMap := make(InstructionMap)
	metadataMap := make(MetadataMap, len(cat.Servers))

	for serviceName, server := range cat.Servers {
		// Skip servers with type "poci"
		if server.Type == "poci" {
			t.logger.Info("skipping server with type 'poci'",
				zap.String("server", serviceName))
			continue
		}

		cfg, req := t.instructionBuilder.BuildServiceConfig(serviceName, server, instructionMap)
		configs[serviceName] = cfg
		requireds[serviceName] = req

		// Extract metadata for this service
		metadataMap[serviceName] = t.extractMetadata(serviceName, server)

		t.logger.Info("transformed server",
			zap.String("server", serviceName),
			zap.Int("properties", len(cfg)),
			zap.Int("required", len(req)))
	}

	t.logger.Info("catalog transformation complete",
		zap.Int("total_servers", len(configs)),
		zap.Int("total_instructions", len(instructionMap)))

	return configs, requireds, instructionMap, metadataMap
}

// extractMetadata extracts metadata from a catalog server
func (t *CatalogTransformer) extractMetadata(serviceName string, server catalog.Server) ServiceMetadata {
	// Use description from catalog
	description := server.Description

	// Construct Docker Hub URL
	// Format: https://hub.docker.com/mcp/server/{simplified-name}/overview
	// Use the service name, but clean it up for the URL
	dockerHubURL := fmt.Sprintf("https://hub.docker.com/mcp/server/%s/overview", serviceName)

	return ServiceMetadata{
		Title:        server.Title,
		Description:  description,
		DockerHubURL: dockerHubURL,
	}
}

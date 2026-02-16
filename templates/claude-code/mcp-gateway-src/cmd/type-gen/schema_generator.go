package main

import (
	"bytes"
	"encoding/json"
	"sort"

	"e2b.dev/mcp-gateway/pkg/naming"
	"e2b.dev/mcp-gateway/pkg/schema"
	"github.com/swaggest/jsonschema-go"
	"go.uber.org/zap"
)

// SchemaGenerator generates JSON schemas from configuration maps
type SchemaGenerator struct {
	logger *zap.Logger
}

// NewSchemaGenerator creates a new SchemaGenerator
func NewSchemaGenerator(logger *zap.Logger) *SchemaGenerator {
	return &SchemaGenerator{logger: logger}
}

// GenerateCombinedSchema generates a combined JSON schema for all services
func (g *SchemaGenerator) GenerateCombinedSchema(configs ConfigMap, requireds RequiredMap, metadataMap MetadataMap) (*jsonschema.Schema, error) {
	combined := &jsonschema.Schema{}
	combined.WithType(jsonschema.Object.Type())
	// Disallow properties not explicitly defined at top level too.
	combined.AdditionalPropertiesEns().WithTypeBoolean(false)

	props := make(map[string]jsonschema.SchemaOrBool, len(configs))
	// Stable order for determinism
	serviceNames := make([]string, 0, len(configs))
	for name := range configs {
		serviceNames = append(serviceNames, name)
	}
	sort.Strings(serviceNames)

	for _, name := range serviceNames {
		// Convert map[string]PropertyInfo to schema.ServiceConfig
		serviceConfig := make(schema.ServiceConfig)
		for k, v := range configs[name] {
			serviceConfig[k] = schema.PropertyInfo{
				Type:        v.Type,
				Description: v.Description,
			}
		}

		serviceSchema := schema.Generate(serviceConfig, requireds[name])
		beautified := naming.BeautifyMcpServerName(name)

		// Add metadata to the service schema
		if metadata, ok := metadataMap[name]; ok {
			if metadata.Title != "" {
				serviceSchema.WithTitle(metadata.Title)
			}
			if metadata.Description != "" {
				serviceSchema.WithDescription(metadata.Description)
			}
			if metadata.DockerHubURL != "" {
				// Add custom x-dockerHubUrl field
				serviceSchema.WithExtraPropertiesItem("x-dockerHubUrl", metadata.DockerHubURL)
			}
		}

		props[beautified] = serviceSchema.ToSchemaOrBool()

		g.logger.Debug("generated schema for service",
			zap.String("server", name),
			zap.String("beautified", beautified))
	}

	combined.WithProperties(props)

	g.logger.Info("combined schema generated",
		zap.Int("services", len(configs)))

	return combined, nil
}

// ToPrettyJSON converts a schema to pretty-printed JSON
func (g *SchemaGenerator) ToPrettyJSON(s *jsonschema.Schema) ([]byte, error) {
	min, err := schema.Minify(s)
	if err != nil {
		return nil, err
	}

	var prettyBuffer bytes.Buffer
	if err := json.Indent(&prettyBuffer, []byte(min), "", "  "); err != nil {
		return nil, err
	}

	return prettyBuffer.Bytes(), nil
}

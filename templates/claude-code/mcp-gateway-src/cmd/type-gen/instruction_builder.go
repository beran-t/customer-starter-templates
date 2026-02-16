package main

import (
	"fmt"
	"sort"
	"strings"

	"e2b.dev/mcp-gateway/pkg/naming"
	"github.com/docker/mcp-gateway/pkg/catalog"
	"go.uber.org/zap"
)

// InstructionBuilder builds instruction mappings from catalog servers
type InstructionBuilder struct {
	logger    *zap.Logger
	extractor *PropertyExtractor
}

// NewInstructionBuilder creates a new InstructionBuilder
func NewInstructionBuilder(logger *zap.Logger) *InstructionBuilder {
	return &InstructionBuilder{
		logger:    logger,
		extractor: NewPropertyExtractor(logger),
	}
}

// BuildServiceConfig constructs the per-service config map and required list
func (b *InstructionBuilder) BuildServiceConfig(
	serviceName string,
	server catalog.Server,
	instructionMap InstructionMap,
) (map[string]PropertyInfo, []string) {
	serviceConfig := make(map[string]PropertyInfo)
	var serviceRequired []string
	beautifiedServiceName := naming.BeautifyMcpServerName(serviceName)

	// Always add a base mapping entry for the service name, even if it has no config
	// This allows the runtime to map beautified names to actual server names
	baseKey := beautifiedServiceName
	instructionMap[baseKey] = Instruction{
		Server: serviceName,
		Type:   "", // Base mapping has no type
	}
	b.logger.Debug("base mapping created",
		zap.String("key", baseKey),
		zap.String("server", serviceName))

	// Extract configuration properties from config array (handles nested configs)
	for configIdx, configItem := range server.Config { // ranging nil slice is safe
		configMap, ok := configItem.(map[string]any)
		if !ok {
			continue
		}
		properties, ok := configMap["properties"].(map[string]any)
		if !ok {
			continue
		}

		// Top-level required for this properties object
		if reqRaw, ok := configMap["required"]; ok {
			switch req := reqRaw.(type) {
			case []any:
				for _, r := range req {
					rk, ok := r.(string)
					if !ok {
						continue
					}
					serviceRequired = append(serviceRequired, naming.ToCamelCase(rk))
				}
			case []string:
				for _, rk := range req {
					serviceRequired = append(serviceRequired, naming.ToCamelCase(rk))
				}
			}
		}

		// Use recursive extraction to handle nested properties
		b.extractor.Extract(serviceName, configIdx, properties, "", []string{}, serviceConfig, &serviceRequired, instructionMap, beautifiedServiceName)
	}

	// Extract secrets and convert to camelCase
	for _, secret := range server.Secrets { // ranging nil slice is safe
		secretKey := naming.ExtractSecretKey(secret.Env, serviceName)
		camelCaseKey := naming.ToCamelCase(strings.ToLower(secretKey))
		serviceConfig[camelCaseKey] = PropertyInfo{Type: "string", Description: ""}
		fullJSONKey := fmt.Sprintf("%s.%s", beautifiedServiceName, camelCaseKey)
		instructionMap[fullJSONKey] = Instruction{
			Server:  serviceName,
			Type:    SecretInstruction,
			EnvName: secret.Env,
		}
		b.logger.Debug("secret instruction created",
			zap.String("key", fullJSONKey),
			zap.String("server", serviceName),
			zap.String("type", "secret"),
			zap.String("envName", secret.Env))
	}

	// If no explicit required provided anywhere for this service, make all keys required
	if len(serviceRequired) == 0 {
		allKeys := make([]string, 0, len(serviceConfig))
		for k := range serviceConfig {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)
		serviceRequired = allKeys
	}

	return serviceConfig, serviceRequired
}

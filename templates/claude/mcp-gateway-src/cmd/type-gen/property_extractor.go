package main

import (
	"fmt"

	"e2b.dev/mcp-gateway/pkg/naming"
	"go.uber.org/zap"
)

// PropertyExtractor handles extraction of nested properties from configuration objects
type PropertyExtractor struct {
	logger *zap.Logger
}

// NewPropertyExtractor creates a new PropertyExtractor
func NewPropertyExtractor(logger *zap.Logger) *PropertyExtractor {
	return &PropertyExtractor{logger: logger}
}

// Extract recursively extracts properties from nested config objects
// and creates flattened keys in camelCase (e.g., "urlConnectionJdbcUrl")
func (e *PropertyExtractor) Extract(
	serviceName string,
	configIdx int,
	properties map[string]any,
	prefix string,
	pathComponents []string,
	result map[string]PropertyInfo,
	required *[]string,
	instructionMap InstructionMap,
	beautifiedServiceName string,
) {
	for propName, propValue := range properties {
		fullKey := propName
		if prefix != "" {
			fullKey = prefix + "_" + propName
		}

		// Build path components for this property
		currentPath := append([]string{}, pathComponents...)
		currentPath = append(currentPath, propName)

		propMap, ok := propValue.(map[string]any)
		if !ok {
			// Simple leaf - default to string
			e.handleSimpleLeaf(serviceName, fullKey, currentPath, result, instructionMap, beautifiedServiceName)
			continue
		}

		// Extract description if present
		description := ""
		if desc, ok := propMap["description"].(string); ok {
			description = desc
		}

		// Nested object with its own properties
		if nestedProps, hasProps := propMap["properties"].(map[string]any); hasProps {
			e.handleNestedObject(serviceName, configIdx, fullKey, propMap, nestedProps, currentPath, result, required, instructionMap, beautifiedServiceName)
			continue
		}

		// Leaf with explicit type information
		if t, ok := propMap["type"].(string); ok {
			e.handleTypedLeaf(serviceName, fullKey, t, propMap, description, currentPath, result, instructionMap, beautifiedServiceName)
			continue
		}

		// Fallback leaf - treat as string
		e.handleFallbackLeaf(serviceName, fullKey, description, currentPath, result, instructionMap, beautifiedServiceName)
	}
}

// handleSimpleLeaf handles a simple leaf property with no type information
func (e *PropertyExtractor) handleSimpleLeaf(
	serviceName string,
	fullKey string,
	currentPath []string,
	result map[string]PropertyInfo,
	instructionMap InstructionMap,
	beautifiedServiceName string,
) {
	camelKey := naming.ToCamelCase(fullKey)
	result[camelKey] = PropertyInfo{Type: "string", Description: ""}
	fullJSONKey := fmt.Sprintf("%s.%s", beautifiedServiceName, camelKey)
	instructionMap[fullJSONKey] = Instruction{
		Server: serviceName,
		Type:   ConfigInstruction,
		Path:   currentPath,
	}
	e.logger.Debug("property extracted",
		zap.String("key", fullJSONKey),
		zap.String("server", serviceName),
		zap.String("type", "config"),
		zap.Strings("path", currentPath),
		zap.String("propType", "string"))
}

// handleNestedObject handles a nested object with its own properties
func (e *PropertyExtractor) handleNestedObject(
	serviceName string,
	configIdx int,
	fullKey string,
	propMap map[string]any,
	nestedProps map[string]any,
	currentPath []string,
	result map[string]PropertyInfo,
	required *[]string,
	instructionMap InstructionMap,
	beautifiedServiceName string,
) {
	// Required keys at this nested level
	switch req := propMap["required"].(type) {
	case []any:
		for _, r := range req {
			rk, ok := r.(string)
			if !ok {
				continue
			}
			*required = append(*required, naming.ToCamelCase(fullKey+"_"+rk))
		}
	case []string:
		for _, rk := range req {
			*required = append(*required, naming.ToCamelCase(fullKey+"_"+rk))
		}
	}

	// Recurse and continue; nothing else to do at this level
	e.Extract(serviceName, configIdx, nestedProps, fullKey, currentPath, result, required, instructionMap, beautifiedServiceName)
}

// handleTypedLeaf handles a leaf property with explicit type information
func (e *PropertyExtractor) handleTypedLeaf(
	serviceName string,
	fullKey string,
	propType string,
	propMap map[string]any,
	description string,
	currentPath []string,
	result map[string]PropertyInfo,
	instructionMap InstructionMap,
	beautifiedServiceName string,
) {
	camelKey := naming.ToCamelCase(fullKey)
	fullJSONKey := fmt.Sprintf("%s.%s", beautifiedServiceName, camelKey)

	// Handle array types with items
	if items, iok := propMap["items"].(map[string]any); propType == "array" && iok {
		if itemType, tok := items["type"].(string); tok {
			result[camelKey] = PropertyInfo{
				Type: map[string]any{
					"type":  "array",
					"items": map[string]any{"type": itemType},
				},
				Description: description,
			}
			instructionMap[fullJSONKey] = Instruction{
				Server: serviceName,
				Type:   ConfigInstruction,
				Path:   currentPath,
			}
			e.logger.Debug("property extracted",
				zap.String("key", fullJSONKey),
				zap.String("server", serviceName),
				zap.String("type", "config"),
				zap.Strings("path", currentPath),
				zap.String("propType", fmt.Sprintf("array of %s", itemType)),
				zap.String("description", description))
			return
		}
	}

	result[camelKey] = PropertyInfo{Type: propType, Description: description}
	instructionMap[fullJSONKey] = Instruction{
		Server: serviceName,
		Type:   ConfigInstruction,
		Path:   currentPath,
	}
	e.logger.Debug("property extracted",
		zap.String("key", fullJSONKey),
		zap.String("server", serviceName),
		zap.String("type", "config"),
		zap.Strings("path", currentPath),
		zap.String("propType", propType),
		zap.String("description", description))
}

// handleFallbackLeaf handles a leaf property with no explicit type, using string as fallback
func (e *PropertyExtractor) handleFallbackLeaf(
	serviceName string,
	fullKey string,
	description string,
	currentPath []string,
	result map[string]PropertyInfo,
	instructionMap InstructionMap,
	beautifiedServiceName string,
) {
	camelKey := naming.ToCamelCase(fullKey)
	result[camelKey] = PropertyInfo{Type: "string", Description: description}
	fullJSONKey := fmt.Sprintf("%s.%s", beautifiedServiceName, camelKey)
	instructionMap[fullJSONKey] = Instruction{
		Server: serviceName,
		Type:   ConfigInstruction,
		Path:   currentPath,
	}
	e.logger.Debug("property extracted (fallback)",
		zap.String("key", fullJSONKey),
		zap.String("server", serviceName),
		zap.String("type", "config"),
		zap.Strings("path", currentPath),
		zap.String("propType", "string"),
		zap.String("description", description))
}

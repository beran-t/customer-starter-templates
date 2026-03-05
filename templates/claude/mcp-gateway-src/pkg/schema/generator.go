package schema

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/swaggest/jsonschema-go"
)

// PropertyInfo holds type and description for a property
type PropertyInfo struct {
	Type        any
	Description string
}

// ServiceConfig represents configuration values for a single service
type ServiceConfig map[string]PropertyInfo

// Generate creates a JSON schema from a ServiceConfig
func Generate(config ServiceConfig, required []string) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.WithType(jsonschema.Object.Type())

	// Build stable, sorted properties map
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	properties := make(map[string]jsonschema.SchemaOrBool, len(keys))
	for _, key := range keys {
		propInfo := config[key]
		prop := &jsonschema.Schema{}

		// Set description if present
		if propInfo.Description != "" {
			prop.WithDescription(propInfo.Description)
		}

		switch v := propInfo.Type.(type) {
		case string:
			// Simple type (string, boolean, number, integer, object, array)
			setSchemaType(prop, v)
		case map[string]any:
			setSchemaType(prop, "string") // fallback
			// Possibly array with items
			if t, ok := v["type"].(string); ok {
				setSchemaType(prop, t)
				if items, iok := v["items"].(map[string]any); t == "array" && iok {
					if itemType, tok := items["type"].(string); tok {
						itemSchema := &jsonschema.Schema{}
						setSchemaType(itemSchema, itemType)
						prop.ItemsEns().WithSchemaOrBool(itemSchema.ToSchemaOrBool())
					}
				}
			}
		default:
			setSchemaType(prop, "string")
		}

		properties[key] = prop.ToSchemaOrBool()
	}

	schema.WithProperties(properties)
	// Disallow properties not explicitly defined.
	schema.AdditionalPropertiesEns().WithTypeBoolean(false)
	if len(required) > 0 {
		schema.WithRequired(required...)
	}
	return schema
}

// setSchemaType maps a JSON Schema type string to swaggest type setter.
func setSchemaType(s *jsonschema.Schema, t string) {
	switch t {
	case "string":
		s.WithType(jsonschema.String.Type())
	case "boolean":
		s.WithType(jsonschema.Boolean.Type())
	case "number":
		s.WithType(jsonschema.Number.Type())
	case "integer":
		s.WithType(jsonschema.Integer.Type())
	case "object":
		s.WithType(jsonschema.Object.Type())
	case "array":
		s.WithType(jsonschema.Array.Type())
	default:
		s.WithType(jsonschema.String.Type())
	}
}

// Minify converts a JSON schema to minified JSON string
func Minify(schema *jsonschema.Schema) (string, error) {
	// Use swaggest's encoder, then compact to ensure minimal output
	raw, err := schema.JSONSchemaBytes()
	if err != nil {
		return "", err
	}

	var compactBuffer bytes.Buffer
	if err := json.Compact(&compactBuffer, raw); err != nil {
		return "", err
	}

	return compactBuffer.String(), nil
}

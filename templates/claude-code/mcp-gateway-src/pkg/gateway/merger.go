package gateway

import (
	"fmt"
	"maps"
	"strings"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/docker/mcp-gateway/pkg/eval"
)

// MergeUserConfigsIntoCatalog merges user configurations into the catalog using instruction map
func MergeUserConfigsIntoCatalog(cat catalog.Catalog, instructionMap InstructionMap, userConfigs map[string]UserConfig) error {
	for beautifiedName, userConfig := range userConfigs {
		// Map beautified name to actual server name in catalog using instruction map
		actualServerName, ok := GetServerNameFromInstructions(instructionMap, beautifiedName)
		if !ok {
			continue
		}

		server, ok := cat.Servers[actualServerName]
		if !ok {
			continue
		}

		// Merge user config into server spec using instructions
		mergedServer, err := mergeUserConfigWithInstructions(actualServerName, server, userConfig, beautifiedName, instructionMap)
		if err != nil {
			return fmt.Errorf("failed to merge config for %s: %w", actualServerName, err)
		}

		// Evaluate placeholders in the merged server
		evaluatedServer, err := evaluatePlaceholders(actualServerName, mergedServer, userConfig, beautifiedName)
		if err != nil {
			return fmt.Errorf("failed to evaluate placeholders for %s: %w", actualServerName, err)
		}

		// Replace in catalog using actual server name
		cat.Servers[actualServerName] = evaluatedServer
	}

	return nil
}

// mergeUserConfigWithInstructions merges user-provided config into the catalog server spec using instructions
func mergeUserConfigWithInstructions(serviceName string, server catalog.Server, userConfig UserConfig, beautifiedName string, instructionMap InstructionMap) (catalog.Server, error) {
	// Create a copy to avoid modifying original
	merged := server

	// Process each user config key
	for userKey, userValue := range userConfig {
		// Build full JSON path: beautifiedName.key
		fullJSONKey := fmt.Sprintf("%s.%s", beautifiedName, userKey)

		instruction, ok := instructionMap[fullJSONKey]
		if !ok {
			continue
		}

		// Verify the instruction is for the correct server
		if instruction.Server != serviceName {
			continue
		}

		// Apply the value to the server based on instruction
		if err := applyInstruction(&merged, instruction, userValue); err != nil {
			return server, fmt.Errorf("failed to apply %s: %w", fullJSONKey, err)
		}
	}

	return merged, nil
}

// applyInstruction applies a value to the server based on the instruction
func applyInstruction(server *catalog.Server, instruction Instruction, value any) error {
	switch instruction.Type {
	case SecretInstruction:
		return applySecret(server, instruction.EnvName, value)
	case ConfigInstruction:
		return applyConfig(server, instruction.Path, value)
	default:
		return fmt.Errorf("unsupported instruction type: %s", instruction.Type)
	}
}

// applySecret sets an environment variable for a secret
func applySecret(server *catalog.Server, envName string, value any) error {
	if envName == "" {
		return fmt.Errorf("secret instruction has no env name")
	}

	// Add or update in server.Env
	valueStr := fmt.Sprintf("%v", value)
	found := false
	for i, env := range server.Env {
		if env.Name == envName {
			server.Env[i].Value = valueStr
			found = true
			break
		}
	}

	if !found {
		server.Env = append(server.Env, catalog.Env{
			Name:  envName,
			Value: valueStr,
		})
	}

	return nil
}

// applyConfig sets a config property value
func applyConfig(server *catalog.Server, path []string, value any) error {
	if len(path) == 0 {
		return fmt.Errorf("config instruction has empty path")
	}

	// Ensure we have at least one config entry
	if len(server.Config) == 0 {
		server.Config = append(server.Config, make(map[string]any))
	}

	// For now, we'll use the first config entry
	// This matches the current behavior where most servers have config[0]
	configMap, ok := server.Config[0].(map[string]any)
	if !ok {
		configMap = make(map[string]any)
		server.Config[0] = configMap
	}

	// Ensure "properties" key exists in config
	if _, exists := configMap["properties"]; !exists {
		configMap["properties"] = make(map[string]any)
	}
	propertiesMap, ok := configMap["properties"].(map[string]any)
	if !ok {
		propertiesMap = make(map[string]any)
		configMap["properties"] = propertiesMap
	}

	// Navigate through the path and set the value
	current := propertiesMap
	for i := range path {
		key := path[i]

		if i == len(path)-1 {
			// Last key - set the value with proper structure
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]any)
			}

			// Ensure it's a map
			propMap, ok := current[key].(map[string]any)
			if !ok {
				propMap = make(map[string]any)
				current[key] = propMap
			}

			// Set the actual value
			propMap["value"] = value
		} else {
			// Intermediate key - need to navigate deeper
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]any)
			}
			nestedMap, ok := current[key].(map[string]any)
			if !ok {
				nestedMap = make(map[string]any)
				current[key] = nestedMap
			}

			// Ensure properties exists for next level
			if _, exists := nestedMap["properties"]; !exists {
				nestedMap["properties"] = make(map[string]any)
			}
			propsMap, ok := nestedMap["properties"].(map[string]any)
			if !ok {
				propsMap = make(map[string]any)
				nestedMap["properties"] = propsMap
			}

			current = propsMap
		}
	}

	return nil
}

// evaluatePlaceholders evaluates placeholder expressions in command, env, volumes, and headers
func evaluatePlaceholders(serviceName string, server catalog.Server, userConfig UserConfig, beautifiedName string) (catalog.Server, error) {
	// Build evaluation context from user config
	// The context needs to be nested: { "dockerhub": { "username": "code42cate", ... } }
	nestedConfig := make(map[string]any)
	maps.Copy(nestedConfig, userConfig)

	// Build catalog-style config - for now just use the user config keys directly
	// This may need adjustment based on how placeholders are used in the catalog
	catalogConfig := make(map[string]any)
	for k, v := range userConfig {
		// Convert camelCase to snake_case for catalog keys
		catalogKey := camelToSnake(k)
		catalogConfig[catalogKey] = v
	}

	evalContext := map[string]any{
		beautifiedName: nestedConfig,  // For user-facing keys like "arxiv"
		serviceName:    catalogConfig, // For catalog placeholders like "arxiv-mcp-server"
	}

	// Evaluate server.remote.Url
	if server.Remote.URL != "" {
		server.Remote.URL = eval.EvaluateList([]string{server.Remote.URL}, evalContext)[0]
	}

	// Evaluate Command
	if len(server.Command) > 0 {
		server.Command = eval.EvaluateList(server.Command, evalContext)
	}

	// Evaluate Volumes
	if len(server.Volumes) > 0 {
		server.Volumes = eval.EvaluateList(server.Volumes, evalContext)
	}

	// Evaluate Env values only if they contain template markers; avoid stripping literal braces
	for i := range server.Env {
		val := server.Env[i].Value
		if strings.Contains(val, "{{") && strings.Contains(val, "}}") {
			result := eval.Evaluate(val, evalContext)
			server.Env[i].Value = fmt.Sprintf("%v", result)
		}
	}

	// Interpolate $VAR and ${VAR} references in env values using current env map (includes secrets)
	if len(server.Env) > 0 {
		envMap := make(map[string]string, len(server.Env))
		for _, e := range server.Env {
			envMap[e.Name] = e.Value
		}
		for i := range server.Env {
			server.Env[i].Value = interpolateEnvVars(server.Env[i].Value, envMap)
		}
	}

	// Evaluate Secret values (if they have example placeholders or similar)
	// Secrets typically reference env vars, so their .Env field shouldn't need evaluation
	// but we can evaluate the Name field if needed
	for i := range server.Secrets {
		result := eval.Evaluate(server.Secrets[i].Name, evalContext)
		server.Secrets[i].Name = fmt.Sprintf("%v", result)
	}

	// Evaluate Remote headers - interpolate ${ENV_VAR} placeholders
	if len(server.Remote.Headers) > 0 {
		server.Remote.Headers = evaluateHeaders(server.Remote.Headers, server.Env)
	}

	return server, nil
}

// evaluateHeaders interpolates ${ENV_VAR} style placeholders in header values
func evaluateHeaders(headers map[string]string, envVars []catalog.Env) map[string]string {
	// Build env lookup map
	envMap := make(map[string]string)
	for _, env := range envVars {
		envMap[env.Name] = env.Value
	}

	evaluated := make(map[string]string, len(headers))
	for key, value := range headers {
		evaluated[key] = interpolateEnvVars(value, envMap)
	}

	return evaluated
}

// interpolateEnvVars replaces ${ENV_VAR} patterns with actual env var values
func interpolateEnvVars(s string, envMap map[string]string) string {
	result := s

	// Replace ${VAR} patterns first
	start := 0
	for {
		idx := strings.Index(result[start:], "${")
		if idx == -1 {
			break
		}
		idx += start

		endIdx := strings.Index(result[idx:], "}")
		if endIdx == -1 {
			break
		}
		endIdx += idx

		// Extract env var name
		envName := result[idx+2 : endIdx]

		// Replace with value if found
		if envValue, ok := envMap[envName]; ok {
			result = result[:idx] + envValue + result[endIdx+1:]
			start = idx + len(envValue)
		} else {
			// Keep original if not found, move past it
			start = endIdx + 1
		}
	}

	// Then replace $VAR patterns (only uppercase letters, digits, and underscore)
	for i := 0; i < len(result); {
		if result[i] == '$' {
			// Skip ${...} which was handled above
			if i+1 < len(result) && result[i+1] == '{' {
				i += 1
				continue
			}
			j := i + 1
			for j < len(result) {
				c := result[j]
				if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
					j++
					continue
				}
				break
			}
			if j > i+1 {
				name := result[i+1 : j]
				if envValue, ok := envMap[name]; ok {
					result = result[:i] + envValue + result[j:]
					i += len(envValue)
					continue
				}
				i = j
				continue
			}
		}
		i++
	}

	return result
}

// camelToSnake converts camelCase to snake_case
func camelToSnake(s string) string {
	var result strings.Builder
	for i, c := range s {
		if i > 0 && c >= 'A' && c <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(c)
	}
	return strings.ToLower(result.String())
}

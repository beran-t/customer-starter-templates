package naming

import "strings"

// ToCamelCase converts snake_case or kebab-case to camelCase
func ToCamelCase(s string) string {
	if strings.Contains(s, "-") {
		s = strings.ReplaceAll(s, "-", "_")
	}

	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return result
}

// BeautifyMcpServerName removes common MCP server name patterns and converts to camelCase
func BeautifyMcpServerName(name string) string {
	name = strings.TrimSuffix(name, "-mcp-server")
	name = strings.TrimSuffix(name, "-mcp")
	name = strings.TrimPrefix(name, "mcp-")
	name = strings.ToLower(name)

	// Convert to camelCase
	name = ToCamelCase(name)

	return name
}

// ExtractSecretKey extracts the key part from an environment variable name
// by removing the service name prefix
// e.g., "GITHUB_TOKEN" with service "github" -> "TOKEN"
func ExtractSecretKey(envName, serviceName string) string {
	prefix := strings.ToUpper(serviceName) + "_"
	return strings.TrimPrefix(envName, prefix)
}

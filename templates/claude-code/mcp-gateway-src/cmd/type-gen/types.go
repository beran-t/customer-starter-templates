package main

// InstructionType defines the type of configuration instruction
type InstructionType string

const (
	SecretInstruction InstructionType = "secret"
	ConfigInstruction InstructionType = "config"
)

// Instruction represents a structured configuration instruction
type Instruction struct {
	Server  string          `json:"server"`            // The catalog server name
	Type    InstructionType `json:"type"`              // "secret" or "config"
	EnvName string          `json:"envName,omitempty"` // For secrets: the environment variable name
	Path    []string        `json:"path,omitempty"`    // For config: the property path
}

// InstructionMap maps JSON keys (including service name) to instructions
// e.g., "github.token" -> Instruction{Server: "github-mcp-server", Type: "secret", EnvName: "GITHUB_TOKEN"}
type InstructionMap map[string]Instruction

// ConfigMap represents the full configuration for all services
// Maps server name to its configuration
type ConfigMap map[string]map[string]PropertyInfo

// RequiredMap collects required property names per service (camelCase, flattened)
type RequiredMap map[string][]string

// PropertyInfo holds type and description for a property
type PropertyInfo struct {
	Type        any
	Description string
}

// ServiceMetadata holds metadata about a service for inclusion in the schema
type ServiceMetadata struct {
	Title        string
	Description  string
	DockerHubURL string
}

// MetadataMap maps service names to their metadata
type MetadataMap map[string]ServiceMetadata

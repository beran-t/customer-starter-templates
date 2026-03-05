package gateway

import (
	"encoding/json"
	"os"
)

// InstructionType defines the type of configuration instruction
type InstructionType string

const (
	SecretInstruction InstructionType = "secret"
	ConfigInstruction InstructionType = "config"
)

// Instruction represents a structured configuration instruction
// that tells the runtime how to apply a user config value to a catalog server
type Instruction struct {
	Server  string          `json:"server"`            // The catalog server name (e.g., "airtable-mcp-server")
	Type    InstructionType `json:"type"`              // "secret" or "config"
	EnvName string          `json:"envName,omitempty"` // For secrets: the environment variable name
	Path    []string        `json:"path,omitempty"`    // For config: the property path (e.g., ["storage_path"])
}

// InstructionMap maps user config keys to instructions
// Key format: "serviceName.propertyName" (e.g., "airtable.apiKey")
type InstructionMap map[string]Instruction

// LoadInstructionMap loads the instruction map from a JSON file
func LoadInstructionMap(path string) (InstructionMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m InstructionMap
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetServerNameFromInstructions gets the actual server name for a beautified name
// by looking it up in the instruction map
func GetServerNameFromInstructions(instructionMap InstructionMap, beautifiedName string) (string, bool) {
	// Look for the base mapping entry (key without dot)
	if instruction, ok := instructionMap[beautifiedName]; ok && instruction.Server != "" {
		return instruction.Server, true
	}

	// Fallback: look for any instruction with this prefix
	prefix := beautifiedName + "."
	for key, instruction := range instructionMap {
		if instruction.Server != "" && len(key) > len(prefix) && key[:len(prefix)] == prefix {
			return instruction.Server, true
		}
	}

	return "", false
}

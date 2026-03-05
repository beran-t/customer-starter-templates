package auth

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

const TokenFilePath = "/etc/mcp-gateway/.token"

// ResolveToken checks for token in order: env var → explicit token → file
func ResolveToken(explicitToken string) (string, error) {
	// 1. Check environment variable
	token := os.Getenv("GATEWAY_ACCESS_TOKEN")
	if token != "" {
		zap.L().Debug("token loaded from environment variable")
		return token, nil
	}

	// 2. Check explicit token (e.g., from CLI flag)
	if explicitToken != "" {
		zap.L().Debug("token loaded from explicit source")
		return explicitToken, nil
	}

	// 3. Check file
	fileToken, err := ReadTokenFromFile()
	if err == nil && fileToken != "" {
		zap.L().Debug("token loaded from file", zap.String("path", TokenFilePath))
		return fileToken, nil
	}

	// No token found (this is okay, auth is optional)
	return "", nil
}

// ReadTokenFromFile reads the token from the token file
func ReadTokenFromFile() (string, error) {
	data, err := os.ReadFile(TokenFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // File doesn't exist, not an error
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// SaveToken writes the token to the token file with 600 permissions
func SaveToken(token string) error {
	// Ensure directory exists
	dir := "/etc/mcp-gateway"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write token with restricted permissions (owner read/write only)
	if err := os.WriteFile(TokenFilePath, []byte(token), 0o600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

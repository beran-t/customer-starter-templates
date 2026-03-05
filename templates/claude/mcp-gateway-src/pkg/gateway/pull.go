package gateway

import (
	"context"
	"fmt"

	"e2b.dev/mcp-gateway/pkg/utils"
	"github.com/docker/mcp-gateway/pkg/catalog"
	"go.uber.org/zap"
)

// PullImages resolves beautified service names to catalog servers and pulls their Docker images.
// Only servers with type "server" and a non-empty Image are pulled. On the first error,
// this function returns immediately with that error.
func PullImages(ctx context.Context, catalogPaths []string, mappingPath string, beautifiedNames []string) error {
	if len(beautifiedNames) == 0 {
		return fmt.Errorf("no services specified to pull")
	}

	// Load instruction map and catalog once
	instructionMap, err := LoadInstructionMap(mappingPath)
	if err != nil {
		return fmt.Errorf("failed to load instruction map: %w", err)
	}

	cat, err := catalog.ReadFrom(ctx, catalogPaths)
	if err != nil {
		return fmt.Errorf("failed to read catalog: %w", err)
	}

	// Create one Docker client for all pulls with custom User-Agent
	dockerClient, err := utils.NewDockerClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer dockerClient.Close()

	for _, beautified := range beautifiedNames {
		// Map beautified name to actual catalog server name
		actualServerName, ok := GetServerNameFromInstructions(instructionMap, beautified)
		if !ok {
			return fmt.Errorf("could not resolve service name %q", beautified)
		}

		server, ok := cat.Servers[actualServerName]
		if !ok {
			return fmt.Errorf("server %q not found in catalog", actualServerName)
		}

		// Only pull images for type "server"
		if server.Type != "server" {
			return fmt.Errorf("service %q is type %q; expected \"server\"", actualServerName, server.Type)
		}

		if server.Image == "" {
			return fmt.Errorf("service %q has no image to pull", actualServerName)
		}

		// Pull the Docker image using shared helper
		if err := utils.PullImage(ctx, server.Image, dockerClient); err != nil {
			return fmt.Errorf("failed to pull image %q: %w", server.Image, err)
		}

		zap.L().Info("Image pulled",
			zap.String("component", "PULL"),
			zap.String("server", actualServerName),
			zap.String("image", server.Image))
	}

	return nil
}

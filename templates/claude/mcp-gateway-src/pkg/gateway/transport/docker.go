package transport

import (
	"context"
	"fmt"
	"os/exec"

	"e2b.dev/mcp-gateway/pkg/utils"
	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// DockerTransport creates MCP sessions by running Docker containers
type DockerTransport struct{}

// CreateSession creates an MCP session by starting a Docker container
func (t *DockerTransport) CreateSession(ctx context.Context, client *mcp.Client, server catalog.Server, serverName string) (*mcp.ClientSession, error) {
	dockerClient, err := utils.NewDockerClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer dockerClient.Close()

	// Pull the Docker image
	if err := utils.PullImage(ctx, server.Image, dockerClient); err != nil {
		zap.L().Error("Failed to pull image",
			zap.String("component", "DOCKER"),
			zap.String("image", server.Image),
			zap.Error(err))
		return nil, fmt.Errorf("failed to pull image %s: %w", server.Image, err)
	}

	// Build docker command arguments
	args := buildDockerArgs(server, serverName)

	// Use context.Background() for long-lived sessions
	commandCtx := ctx
	if server.LongLived {
		commandCtx = context.Background()
	}

	// Connect to Docker container
	session, err := client.Connect(ctx, &mcp.CommandTransport{
		Command: exec.CommandContext(commandCtx, "docker", args...),
	}, nil)
	if err != nil {
		zap.L().Error("Failed to connect to Docker container",
			zap.String("component", "DOCKER"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Docker container: %w", err)
	}

	return session, nil
}

// buildDockerArgs constructs the docker run command arguments
func buildDockerArgs(server catalog.Server, serverName string) []string {
	args := []string{"run"}

	// Base security and resource settings
	args = append(args, "--rm", "-i", "--init", "--security-opt", "no-new-privileges")
	args = append(args, "--cpus", "1")
	args = append(args, "--memory", "1g")
	args = append(args, "--pull", "never")

	// Docker MCP labels
	args = append(args,
		"-l", "docker-mcp=true",
		"-l", "docker-mcp-tool-type=mcp",
		"-l", "docker-mcp-name="+serverName,
		"-l", "docker-mcp-transport=stdio",
	)

	// Network isolation
	if server.DisableNetwork {
		args = append(args, "--network", "none")
	}

	// Volumes from catalog (already evaluated with placeholders)
	for _, volume := range server.Volumes {
		if volume != "" {
			args = append(args, "-v", volume)
		}
	}

	// User from catalog
	if server.User != "" {
		args = append(args, "-u", server.User)
	}

	// Environment variables from catalog (already merged with user config)
	for _, e := range server.Env {
		if e.Name != "" && e.Value != "" {
			args = append(args, "-e", fmt.Sprintf("%s=%s", e.Name, e.Value))
		}
	}

	// Image and command
	args = append(args, server.Image)
	args = append(args, server.Command...)

	return args
}

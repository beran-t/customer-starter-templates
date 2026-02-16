package transport

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// GitHubTransport creates MCP sessions by cloning GitHub repos and running stdio commands
type GitHubTransport struct{}

// CreateSession creates an MCP session by cloning a GitHub repo and starting an MCP server
func (t *GitHubTransport) CreateSession(ctx context.Context, client *mcp.Client, server catalog.Server, serverName string) (*mcp.ClientSession, error) {
	// Parse server name to extract username/repo from github/{username}/{repo}
	parts := strings.Split(serverName, "/")
	if len(parts) < 3 || parts[0] != "github" {
		return nil, fmt.Errorf("invalid github server name format: %s (expected github/{username}/{repo})", serverName)
	}
	username := parts[1]
	repo := parts[2]

	// Generate clone path
	clonePath := filepath.Join("/var/lib/mcp-gateway/github", username, repo)

	// Clean up any existing directory
	if err := os.RemoveAll(clonePath); err != nil {
		zap.L().Warn("Failed to clean existing directory",
			zap.String("component", "GITHUB"),
			zap.String("path", clonePath),
			zap.Error(err))
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(clonePath), 0o755); err != nil {
		zap.L().Error("Failed to create parent directory",
			zap.String("component", "GITHUB"),
			zap.String("path", filepath.Dir(clonePath)),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Clone the repository
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", username, repo)
	zap.L().Info("Cloning GitHub repository",
		zap.String("component", "GITHUB"),
		zap.String("url", repoURL),
		zap.String("path", clonePath))

	cloneCmd := exec.CommandContext(ctx, "git", "clone", repoURL, clonePath)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		zap.L().Error("Failed to clone repository",
			zap.String("component", "GITHUB"),
			zap.String("url", repoURL),
			zap.String("output", string(output)),
			zap.Error(err))
		os.RemoveAll(clonePath) // Clean up on failure
		return nil, fmt.Errorf("failed to clone repository %s: %w", repoURL, err)
	}

	// Extract install command from server.Env (special env var hack)
	var installCommand string
	for _, env := range server.Env {
		if env.Name == "INSTALL_COMMAND" {
			installCommand = env.Value
			break
		}
	}

	// Execute install command if provided
	if installCommand != "" {
		zap.L().Info("Running install command",
			zap.String("component", "GITHUB"),
			zap.String("command", installCommand),
			zap.String("workdir", clonePath))

		installCmd := exec.CommandContext(ctx, "sh", "-c", installCommand)
		installCmd.Dir = clonePath
		if output, err := installCmd.CombinedOutput(); err != nil {
			zap.L().Error("Failed to run install command",
				zap.String("component", "GITHUB"),
				zap.String("command", installCommand),
				zap.String("output", string(output)),
				zap.Error(err))
			os.RemoveAll(clonePath) // Clean up on failure
			return nil, fmt.Errorf("failed to run install command: %w", err)
		}
	}

	// Use context.Background() for long-lived sessions
	commandCtx := ctx
	if server.LongLived {
		commandCtx = context.Background()
	}

	// Execute run command with CommandTransport
	zap.L().Info("Starting MCP server",
		zap.String("component", "GITHUB"),
		zap.Strings("command", server.Command),
		zap.String("workdir", clonePath))

	if len(server.Command) == 0 {
		return nil, fmt.Errorf("empty command for GitHub server")
	}

	commandStr := strings.Join(server.Command, " ")
	zap.L().Info("About to execute command",
		zap.String("component", "GITHUB"),
		zap.String("commandStr", commandStr),
		zap.String("workdir", clonePath))

	cmdparts := strings.Fields(commandStr)
	if len(cmdparts) == 0 {
		return nil, fmt.Errorf("empty command after parsing")
	}

	var runCmd *exec.Cmd
	if len(cmdparts) == 1 {
		runCmd = exec.CommandContext(commandCtx, cmdparts[0])
	} else {
		runCmd = exec.CommandContext(commandCtx, cmdparts[0], cmdparts[1:]...)
	}
	runCmd.Dir = clonePath
	runCmd.Env = os.Environ()

	// Connect to the command's stdio
	session, err := client.Connect(ctx, &mcp.CommandTransport{
		Command:           runCmd,
		TerminateDuration: 30 * time.Second,
	}, nil)
	if err != nil {
		zap.L().Error("Failed to connect to GitHub MCP server",
			zap.String("component", "GITHUB"),
			zap.String("server", serverName),
			zap.Error(err))
		os.RemoveAll(clonePath) // Clean up on failure
		return nil, fmt.Errorf("failed to connect to GitHub MCP server: %w", err)
	}

	return session, nil
}

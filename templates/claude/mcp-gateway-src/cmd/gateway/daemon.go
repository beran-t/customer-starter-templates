package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"e2b.dev/mcp-gateway/pkg/auth"
	"e2b.dev/mcp-gateway/pkg/gateway"
	"e2b.dev/mcp-gateway/pkg/health"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// runLauncher spawns the daemon process, checks health, and exits
func runLauncher(c *cli.Context, host string, port int) error {
	args := []string{}

	for _, arg := range os.Args[1:] {
		if arg == "--daemon" {
			continue
		}
		args = append(args, arg)
	}
	args = append(args, "--daemon")

	// Get the executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	logDir := "/var/log/mcp-gateway"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		logDir = os.TempDir()
	}

	logFile, err := os.OpenFile(fmt.Sprintf("%s/gateway.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	cmd := exec.Command(executable, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon process: %w", err)
	}

	zap.L().Info("daemon process started",
		zap.Int("pid", cmd.Process.Pid),
		zap.String("log_file", fmt.Sprintf("%s/gateway.log", logDir)))

	if err := health.WaitForAlive(c.Context, fmt.Sprintf("http://%s:%d/health", host, port), 100*time.Millisecond); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	zap.L().Info("gateway is healthy and running in background")
	return nil
}

// runDaemon runs the actual server process
func runDaemon(ctx context.Context, c *cli.Context, host string, port int) error {
	catalogs := c.StringSlice("catalog")
	mapping := c.String("mapping")
	configStr := c.String("config")

	// Resolve token: env var → CLI flag → file
	token, err := auth.ResolveToken(c.String("token"))
	if err != nil {
		return fmt.Errorf("failed to resolve token: %w", err)
	}

	// Save token to file if requested
	if c.Bool("save-token") && token != "" {
		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		zap.L().Info("token saved", zap.String("path", auth.TokenFilePath))
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	g, err := gateway.New(
		ctx,
		catalogs,
		mapping,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize gateway: %w", err)
	}

	// Load config if provided
	if configStr != "" {
		if err := g.LoadConfig(ctx, []byte(configStr)); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		zap.L().Info("configuration loaded")
	}

	// Setup and run HTTP server
	srv := gateway.SetupHTTPServer(g.Server(), addr, token)

	zap.L().Info("server starting", zap.String("addr", addr))
	gateway.Run(srv, g.Pool())

	return nil
}

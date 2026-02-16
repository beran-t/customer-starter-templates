package main

import (
	"context"
	"fmt"
	"os"

	"e2b.dev/mcp-gateway/pkg/gateway"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	app := &cli.App{
		Name:  "mcp-gateway",
		Usage: "MCP Gateway Server",
		Commands: []*cli.Command{
			{
				Name:  "pull",
				Usage: "Pull Docker images for specified services (by beautified name)",
				Action: func(c *cli.Context) error {
					ctx := context.Background()

					// Ensure Docker config path like in run
					os.Setenv("DOCKER_CONFIG", "/root/.docker")

					args := c.Args().Slice()
					if len(args) == 0 {
						return fmt.Errorf("please specify at least one service to pull")
					}

					// Delegate to gateway logic
					return gateway.PullImages(ctx, c.StringSlice("catalog"), c.String("mapping"), args)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Value:   "0.0.0.0",
				Usage:   "server host",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   50005,
				Usage:   "server port",
			},
			&cli.StringSliceFlag{
				Name:    "catalog",
				Aliases: []string{"c"},
				Value:   cli.NewStringSlice("/etc/mcp-gateway/docker-catalog.yaml"),
				Usage:   "catalog file(s) to load (can be specified multiple times)",
			},
			&cli.StringFlag{
				Name:    "mapping",
				Aliases: []string{"m"},
				Value:   "/etc/mcp-gateway/mapping.json",
				Usage:   "mapping file to load",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "configuration JSON",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "authentication token (enables auth middleware)",
			},
			&cli.BoolFlag{
				Name:  "save-token",
				Value: true,
				Usage: "save the token to file for future use",
			},
			&cli.BoolFlag{
				Name:    "foreground",
				Aliases: []string{"f"},
				Value:   false,
				Usage:   "run in foreground instead of as a daemon",
			},
			&cli.BoolFlag{
				Name:   "daemon",
				Usage:  "internal flag: run as daemon (used internally)",
				Hidden: true,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		zap.L().Fatal("application error", zap.Error(err))
	}
}

func run(c *cli.Context) error {
	ctx := context.Background()

	os.Setenv("DOCKER_CONFIG", "/root/.docker")

	host := c.String("host")
	port := c.Int("port")
	isDaemon := c.Bool("daemon")
	isForeground := c.Bool("foreground")

	// If running in foreground mode, start server directly
	if isForeground {
		return runDaemon(ctx, c, host, port)
	}

	// If not running as daemon, spawn background process and check health
	if !isDaemon {
		return runLauncher(c, host, port)
	}

	// Running as daemon, start the server
	return runDaemon(ctx, c, host, port)
}

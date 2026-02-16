package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// WaitForAlive polls a health endpoint until it returns 200 OK
func WaitForAlive(ctx context.Context, healthURL string, pollInterval time.Duration) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	zap.L().Info("waiting for health check", zap.String("url", healthURL))

	// Try immediately first
	if err := checkHealth(ctx, client, healthURL); err == nil {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("health check cancelled: %w", ctx.Err())
		case <-ticker.C:
			if err := checkHealth(ctx, client, healthURL); err == nil {
				return nil
			}
		}
	}
}

func checkHealth(ctx context.Context, client *http.Client, healthURL string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.L().Debug("health check failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		zap.L().Info("health check passed")
		return nil
	}

	zap.L().Debug("health check returned non-200", zap.Int("status", resp.StatusCode))
	return fmt.Errorf("health check returned status %d", resp.StatusCode)
}

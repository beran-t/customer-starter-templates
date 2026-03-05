package utils

import (
	"context"

	apiimage "github.com/docker/docker/api/types/image"
	client2 "github.com/docker/docker/client"
	dockerclient "github.com/docker/go-sdk/client"
	"github.com/docker/go-sdk/image"
)

const UserAgent = "E2B/0.0.1"

// NewDockerClient creates a Docker client with the provided User-Agent header.
func NewDockerClient(ctx context.Context) (*dockerclient.Client, error) {
	return dockerclient.New(ctx, dockerclient.FromDockerOpt(
		client2.WithUserAgent(UserAgent),
	))
}

// PullImage pulls the specified Docker image using the provided client.
func PullImage(ctx context.Context, serverImage string, client *dockerclient.Client) error {
	return image.Pull(ctx,
		serverImage,
		image.WithPullClient(client),
		image.WithPullOptions(apiimage.PullOptions{}),
	)
}

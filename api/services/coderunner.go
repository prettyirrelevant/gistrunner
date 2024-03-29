package services

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"go.uber.org/zap"

	"github.com/prettyirrelevant/gistrunner/helpers"
	"github.com/prettyirrelevant/gistrunner/types/languages"
)

type CodeRunner struct {
	url    string
	logger *zap.Logger
	client *client.Client
}

func NewCodeRunner(url string, logger *zap.Logger) (*CodeRunner, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost(url),
		client.WithScheme("https"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	runner := &CodeRunner{url: url, client: cli, logger: logger}
	if err = runner.clearCache(); err != nil {
		return nil, err
	}
	if err = runner.setupCodeRunnerImages(); err != nil {
		return nil, err
	}

	return runner, nil
}

func (c *CodeRunner) setupCodeRunnerImages() error {
	for _, image := range languages.SupportedLanguagesImages {
		c.logger.Debug("pulling docker image from dockerhub", zap.String("image", image))
		ctx := context.Background()
		if _, err := c.client.ImagePull(ctx, image, types.ImagePullOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRunner) Run(lang languages.ProgrammingLanguage, code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	resp, err := c.client.ContainerCreate(
		ctx,
		&container.Config{
			Tty:   true,
			Cmd:   lang.ContainerCommand(),
			Image: languages.SupportedLanguagesImages[lang],
			User:  "nobody", // Run as a non-root user
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory: 64 * 1024 * 1024, // 64MB memory limit
			},
			CapDrop:        []string{"ALL"},                 // Drop all capabilities
			UsernsMode:     container.UsernsMode("private"), // Enable user namespaces
			Privileged:     false,                           // Disable privileged mode
			ReadonlyRootfs: true,                            // Read-only filesystem
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					Source:   "/tmp",
					Target:   "/tmp",
					ReadOnly: false,
				},
			},
		},
		nil, nil, "",
	)
	if err != nil {
		return "", err
	}
	defer c.closeContainer(resp.ID)

	tarReader, err := helpers.CreateTarArchiveForSourceCode(code, languages.SupportedLanguages[lang])
	if err != nil {
		return "", err
	}

	if err = c.client.CopyToContainer(ctx, resp.ID, "/tmp", tarReader, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
		return "", err
	}

	if err := c.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := c.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	rawLogs, err := c.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Details: true})
	if err != nil {
		return "", nil
	}
	defer rawLogs.Close()

	logs, err := io.ReadAll(rawLogs)
	if err != nil {
		return "", nil
	}

	return string(logs), nil
}

func (c *CodeRunner) closeContainer(id string) error {
	return c.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (c *CodeRunner) clearCache() error {
	c.logger.Debug("pruning build cache")
	if _, err := c.client.BuildCachePrune(context.Background(), types.BuildCachePruneOptions{All: true}); err != nil {
		return err
	}

	c.logger.Debug("pruning containers cache")
	if _, err := c.client.ContainersPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}

	return nil
}

func (c *CodeRunner) Close() error {
	return c.client.Close()
}

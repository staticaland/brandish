package markdownlint

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

func base(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("alpine:latest")
}

func Version(client *dagger.Client) string {
	ctx := context.Background()

	out, err := base(client).
		WithExec([]string{"cat", "/etc/alpine-release"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

func baseMarkdownlint(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("davidanson/markdownlint-cli2:latest").
		Pipeline("markdownlint")
}

type MarkdownlintOption func(*MarkdownlintConfig)

type MarkdownlintConfig struct {
	Config string
	Globs  string
	Fix    bool
}

func WithGlobs(globs string) MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Globs = globs
	}
}

func WithFix() MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Fix = true
	}
}

func WithConfig(config string) MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Config = config
	}
}

func Lint(client *dagger.Client, opts ...MarkdownlintOption) string {
	ctx := context.Background()

	// Default configuration
	cfg := &MarkdownlintConfig{
		Config: "",
		Globs:  "*.{md,markdown}",
		Fix:    false,
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// We will not be using the default ENTRYPOINT and CMD of the image
	markdownlintCommand := []string{
		"/usr/local/bin/markdownlint-cli2",
		cfg.Globs,
	}
	if cfg.Config != "" {
		markdownlintCommand = append(markdownlintCommand, "--config", cfg.Config)
	}
	if cfg.Fix {
		markdownlintCommand = append(markdownlintCommand, "--fix")
	}

	// Always return true so that the container exits with 0
	markdownlintCommand = append(markdownlintCommand, ";", "true")

	markdownlintCommandStr := strings.Join(markdownlintCommand, " ")

	args := []string{"-c", markdownlintCommandStr}

	container, err := baseMarkdownlint(client).
		WithDirectory("/home/node", client.Host().Directory("."), dagger.ContainerWithDirectoryOpts{
			Include: []string{
				cfg.Globs,
			},
			Owner: "node",
		}).
		WithWorkdir("/home/node").
		WithEntrypoint([]string{"sh"}).
		WithExec(args).
		Sync(ctx)

	if err != nil {
		// Unexpected error, could be network failure.
		fmt.Println(err)
	}

	out, err := container.Stdout(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// If Fix is true, export the changes back to the host
	if cfg.Fix {
		_, err := container.Directory(".").Export(ctx, ".")
		if err != nil {
			fmt.Println(err)
		}
	}

	return out
}

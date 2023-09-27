package pipelines

import (
	"context"

	"dagger.io/dagger"
)

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

func WithCommand(config string) MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Config = config
	}
}

// create base image
func base(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("alpine:latest")
}

func baseMarkdownlint(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("davidanson/markdownlint-cli2")
}

// run command in base image
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

func Markdownlint(client *dagger.Client, opts ...MarkdownlintOption) string {
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

	args := []string{cfg.Globs}
	if cfg.Config != "" {
		args = append(args, "--config", cfg.Config)
	}
	if cfg.Fix {
		args = append(args, "--fix")
	}

	src := client.Host().Directory(".")

	out, err := baseMarkdownlint(client).
		WithDirectory("/src", src).WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

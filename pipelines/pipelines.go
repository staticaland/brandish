package pipelines

import (
	"context"

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
		From("davidanson/markdownlint-cli2:latest")
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

func WithFix(fix bool) MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Fix = fix
	}
}

func WithConfig(config string) MarkdownlintOption {
	return func(cfg *MarkdownlintConfig) {
		cfg.Config = config
	}
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

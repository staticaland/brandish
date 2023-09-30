package terraform

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func baseTerraform(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("hashicorp/terraform:latest").
		Pipeline("terraform")
}

type TerraformOptions func(*TerraformConfig)

type TerraformConfig struct {
	Path      string
	Recursive bool
}

func WithPath(path string) TerraformOptions {
	return func(cfg *TerraformConfig) {
		cfg.Path = path
	}
}

func WithRecursive(recursive bool) TerraformOptions {
	return func(cfg *TerraformConfig) {
		cfg.Recursive = recursive
	}
}

func Fmt(client *dagger.Client, opts ...TerraformOptions) string {
	ctx := context.Background()

	// Default configuration
	cfg := &TerraformConfig{
		Path:      ".",
		Recursive: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	args := []string{
		"fmt",
	}

	if cfg.Recursive {
		args = append(args, "-recursive")
	}

	if cfg.Path != "." {
		args = append(args, cfg.Path)
	}

	container, err := baseTerraform(client).
		Pipeline("fmt").
		WithDirectory("/workdir", client.Host().Directory("."), dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"**/*.tf",
			},
		}).
		WithWorkdir("/workdir").
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

	// Export the changes back to the host
	_, err = container.Directory(".").Export(ctx, ".")
	if err != nil {
		fmt.Println(err)
	}

	return out
}

func Plan(client *dagger.Client) string {
	ctx := context.Background()

	args := []string{
		"plan",
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	container, err := baseTerraform(client).
		Pipeline("plan").
		WithEnvVariable("AWS_PROFILE", os.Getenv("AWS_PROFILE")).
		WithDirectory("/root/.aws", client.Host().Directory(fmt.Sprintf("%s/.aws", homeDir))).
		WithDirectory("/workdir", client.Host().Directory("."), dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"**/*.tf",
			},
		}).
		WithWorkdir("/workdir").
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

	return out
}

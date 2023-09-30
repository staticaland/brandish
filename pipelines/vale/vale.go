package vale

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func baseVale(client *dagger.Client) *dagger.Container {
	return client.
		Container().
		From("jdkato/vale:latest").
		Pipeline("vale")
}

type ValeOption func(*ValeConfig)

type ValeConfig struct {
	HostDir    string   // HostDir is the directory on the host machine that will be mounted into the container at /workdir.
	Paths      []string // Used for input in vale [options] [input...]
	Globs      string   // Used for --glob. A glob pattern (--glob='*.{md,txt}.').
	StylesPath string   // Used for for StylesPath inside the config. So it should be the value like in StylesPath = a/path/to/your/style
	ConfigPath string   // Used for --config. A file path (--config='some/file/path/.vale.ini').
}

func WithHostDir(hostDir string) ValeOption {
	return func(cfg *ValeConfig) {
		cfg.HostDir = hostDir
	}
}

func WithPaths(paths []string) ValeOption {
	return func(cfg *ValeConfig) {
		cfg.Paths = paths
	}
}

func WithGlobs(globs string) ValeOption {
	return func(cfg *ValeConfig) {
		cfg.Globs = globs
	}
}

func WithStylesPath(stylesPath string) ValeOption {
	return func(cfg *ValeConfig) {
		cfg.StylesPath = stylesPath
	}
}

func WithConfigPath(configPath string) ValeOption {
	return func(cfg *ValeConfig) {
		cfg.ConfigPath = configPath
	}
}

func BuildArgs(cfg *ValeConfig, opts ...ValeOption) []string {

	for _, opt := range opts {
		opt(cfg)
	}

	args := []string{}

	if cfg.Globs != "" {
		args = append(args, fmt.Sprintf("--glob='%s'", cfg.Globs))
	}

	if cfg.ConfigPath != ".vale.ini" {
		args = append(args, fmt.Sprintf("--config='%s'", cfg.ConfigPath))
	}

	if cfg.Paths != nil {
		args = append(args, cfg.Paths...)
	}

	// Example of final command when joining args:
	// vale --glob='*.{md,txt}' --config='some/file/path/.vale.ini' some/file/path
	return args
}

func Vale(client *dagger.Client, opts ...ValeOption) string {
	ctx := context.Background()

	// Default configuration
	cfg := &ValeConfig{
		HostDir:    ".",
		Paths:      []string{"."},
		Globs:      "",
		StylesPath: "styles",
		ConfigPath: ".vale.ini",
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	args := BuildArgs(cfg, opts...)

	container, err := baseVale(client).
		WithDirectory("/workdir", client.Host().Directory(cfg.HostDir), dagger.ContainerWithDirectoryOpts{
			Include: append([]string{
				cfg.ConfigPath,
				cfg.StylesPath,
				"**/*.md",
			}, cfg.Paths...),
		}).
		WithWorkdir("/workdir").
		WithExec(args).
		Sync(ctx)

	if err != nil {
		fmt.Println(err)
	}

	out, err := container.Stdout(ctx)
	if err != nil {
		fmt.Println(err)
	}

	return out
}

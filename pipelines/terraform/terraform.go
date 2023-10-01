package terraform

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

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

type Terraform struct {
	ctx       context.Context
	client    *dagger.Client
	container *dagger.Container
	config    *TerraformConfig
}

func New(ctx context.Context, client *dagger.Client) *Terraform {
	tf := &Terraform{ctx: ctx, client: client}
	tf.container = tf.NewContainer()
	tf.config = &TerraformConfig{
		Path:      ".",
		Recursive: true,
	}
	return tf
}

func (tf *Terraform) NewContainer() *dagger.Container {
	return tf.client.
		Container().
		From("hashicorp/terraform:latest").
		Pipeline("terraform")
}

func (tf *Terraform) WithTerraformFiles() *Terraform {
	terraformFiles := tf.client.Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{
			"**/*.tf",
		},
	})

	tf.container = tf.container.
		WithDirectory("/workdir", terraformFiles).
		WithWorkdir("/workdir")

	return tf
}

func (tf *Terraform) WithExecuteFmt(opts ...TerraformOptions) (*dagger.Container, error) {
	// Apply options
	for _, opt := range opts {
		opt(tf.config)
	}

	args := []string{
		"fmt",
	}

	if tf.config.Recursive {
		args = append(args, "-recursive")
	}

	if tf.config.Path != "." {
		args = append(args, tf.config.Path)
	}

	tf = tf.WithTerraformFiles()

	container, err := tf.container.
		Pipeline("fmt").
		WithExec(args).
		Sync(tf.ctx)

	return container, err
}

func (tf *Terraform) Fmt(opts ...TerraformOptions) string {
	container, err := tf.WithExecuteFmt(opts...)
	if err != nil {
		// Unexpected error, could be network failure.
		fmt.Println(err)
		return ""
	}

	out, err := container.Stdout(tf.ctx)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Export the changes back to the host
	_, err = container.Directory(".").Export(tf.ctx, ".")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return out
}

func (tf *Terraform) WithAWSAuth() *Terraform {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return tf
	}

	tf.container = tf.container.
		WithEnvVariable("AWS_PROFILE", os.Getenv("AWS_PROFILE")).
		WithDirectory("/root/.aws", tf.client.Host().Directory(fmt.Sprintf("%s/.aws", homeDir)))

	return tf
}

func (tf *Terraform) WithExecutePlan() (*dagger.Container, error) {
	args := []string{
		"plan",
	}

	container, err := tf.container.
		Pipeline("plan").
		WithExec(args).
		Sync(tf.ctx)

	return container, err
}

func (tf *Terraform) Plan() string {
	tf = tf.WithTerraformFiles().WithAWSAuth()

	container, err := tf.WithExecutePlan()
	if err != nil {
		// Unexpected error, could be network failure.
		fmt.Println(err)
		return ""
	}

	out, err := container.Stdout(tf.ctx)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return out
}

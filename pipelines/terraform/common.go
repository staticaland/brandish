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

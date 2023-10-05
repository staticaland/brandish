package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

type CommonConfig struct {
	Path string
}

type FmtConfig struct {
	Recursive bool
}

type PlanConfig struct {
	// Add any specific configuration for the Plan command here
}

type TerraformOptions func(*Terraform)

type Terraform struct {
	ctx          context.Context
	client       *dagger.Client
	container    *dagger.Container
	commonConfig *CommonConfig
	fmtConfig    *FmtConfig
	planConfig   *PlanConfig
}

func WithCommonPath(path string) TerraformOptions {
	return func(tf *Terraform) {
		tf.commonConfig.Path = path
	}
}

func WithFmtRecursive(recursive bool) TerraformOptions {
	return func(tf *Terraform) {
		tf.fmtConfig.Recursive = recursive
	}
}

func New(ctx context.Context, client *dagger.Client, opts ...TerraformOptions) *Terraform {
	tf := &Terraform{ctx: ctx, client: client}
	tf.container = tf.NewContainer()
	tf.commonConfig = &CommonConfig{
		Path: ".",
	}
	tf.fmtConfig = &FmtConfig{
		Recursive: true,
	}
	tf.planConfig = &PlanConfig{
		// Initialize any PlanConfig fields here
	}
	for _, opt := range opts {
		opt(tf)
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

	workdir := filepath.Join("/workdir", tf.commonConfig.Path)

	terraformFiles := tf.client.Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{
			"**/*.tf",
		},
	})

	tf.container = tf.container.
		WithDirectory("/workdir", terraformFiles).
		WithWorkdir(workdir)

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

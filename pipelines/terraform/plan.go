package terraform

import (
	"dagger.io/dagger"
)

// TODO: I like having a output struct, but how do I handle errors and exit codes?
// There is something about https://pkg.go.dev/dagger.io/dagger#ExecError

type PlanOutputs struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (tf *Terraform) WithExecInit() *Terraform {
	args := []string{
		"init",
	}

	tf.container = tf.container.
		WithExec(args)

	return tf
}

func (tf *Terraform) WithExecPlan() (*dagger.Container, error) {
	args := []string{
		"plan",
	}

	container, err := tf.container.
		Pipeline("plan").
		WithExec(args).
		Sync(tf.ctx)

	return container, err
}

func (tf *Terraform) Plan() (PlanOutputs, error) {
	tf = tf.WithTerraformFiles().WithAWSAuth().WithExecInit()

	container, err := tf.WithExecPlan()
	if err != nil {
		// Unexpected error, could be network failure.
		return PlanOutputs{}, err
	}

	out, err := container.Stdout(tf.ctx)
	if err != nil {
		return PlanOutputs{}, err
	}

	errOut, err := container.Stderr(tf.ctx)
	if err != nil {
		return PlanOutputs{}, err
	}

	return PlanOutputs{
		Stdout:   out,
		Stderr:   errOut,
		ExitCode: 0, // TODO: Replace with the actual exit code
	}, nil

}

package terraform

import (
	"fmt"

	"dagger.io/dagger"
)

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

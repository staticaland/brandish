package terraform

import (
	"fmt"

	"dagger.io/dagger"
)

func (tf *Terraform) WithExecuteFmt(opts ...TerraformOptions) (*dagger.Container, error) {
	// Apply options
	for _, opt := range opts {
		opt(tf)
	}

	args := []string{
		"fmt",
	}

	if tf.fmtConfig.Recursive {
		args = append(args, "-recursive")
	}

	if tf.commonConfig.Path != "." {
		args = append(args, tf.commonConfig.Path)
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

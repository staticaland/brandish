package pipelines

import (
	"context"

	"dagger.io/dagger"
)

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

func Markdownlint(client *dagger.Client, path string) string {
	ctx := context.Background()

	out, err := baseMarkdownlint(client).
		WithExec([]string{"markdownlint", path}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

//go:build mage

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/sh"
	"github.com/staticaland/brandish/pipelines"
)

// Build the project.
func Build() error {
	return sh.Run("go", "build", "-o", "main", "main.go")
}

// Run the project.
func Run() error {
	return sh.Run("go", "run", "main.go")
}

// Lint
func Lint(ctx context.Context) error {

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(pipelines.Markdownlint(client))

	return nil
}

// Lint and fix
func LintFix(ctx context.Context) error {

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(
		pipelines.Markdownlint(
			client,
			pipelines.WithGlobs("README.md"),
			pipelines.WithFix(),
		),
	)

	return nil
}

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/staticaland/brandish/pipelines/markdownlint"
	"github.com/staticaland/brandish/pipelines/terraform"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(
		markdownlint.Lint(
			client,
			markdownlint.WithGlobs("README.md"),
			markdownlint.WithFix(),
		),
	)
	fmt.Println(terraform.Fmt(client))
}

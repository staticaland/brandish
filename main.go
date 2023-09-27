package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/staticaland/universe/pipelines"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(pipelines.Markdownlint(client))
	fmt.Println(
		pipelines.Markdownlint(
			client,
			pipelines.WithGlobs("README.md"),
			pipelines.WithFix(),
		),
	)
}

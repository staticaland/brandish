package main

import (
	"context"
	"fmt"
	"os"

	d "dagger.io/dagger"
	md "github.com/staticaland/brandish/pipelines/markdownlint"
	tf "github.com/staticaland/brandish/pipelines/terraform"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := d.Connect(ctx, d.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(
		md.Lint(
			client,
			md.WithGlobs("README.md"),
			md.WithFix(),
		),
	)
	fmt.Println(tf.Fmt(client))
}

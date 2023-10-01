package main

import (
	"context"
	"fmt"
	"os"

	d "dagger.io/dagger"
	md "github.com/staticaland/brandish/pipelines/markdownlint"
	terraform "github.com/staticaland/brandish/pipelines/terraform"
	vale "github.com/staticaland/brandish/pipelines/vale"
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
	tf := terraform.New(ctx, client)

	fmt.Println(tf.Fmt())
	fmt.Println(tf.Plan())

	fmt.Println(vale.Vale(
		client,
		vale.WithHostDir("workdirs/vale")),
	)

}

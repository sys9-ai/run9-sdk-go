package run9_test

import (
	"context"

	run9 "github.com/sys9-ai/run9-sdk-go"
)

func ExampleNewClient() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	projectClient := client.WithProject("default")
	_, _ = client, projectClient
}

func ExampleClient_WithProject() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	_, _ = project.ListBoxes(context.Background(), run9.ListBoxesRequest{})
}

func ExampleClient_StartExecStream() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	stream, err := project.StartExecStream(context.Background(), "devbox", run9.ExecRequest{
		Command: []string{"/bin/sh", "-lc", "echo hello"},
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	_, _ = stream.ReadEvent()
}

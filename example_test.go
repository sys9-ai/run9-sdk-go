package run9_test

import (
	"context"
	"io"
	"time"

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

func ExampleClient_RunExec() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	_, err = project.RunExec(context.Background(), "devbox", run9.ExecRequest{
		Command: []string{"/bin/sh", "-lc", "echo hello"},
	}, run9.ExecOutputWriters{
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	if err != nil {
		panic(err)
	}

	_ = project
}

func ExampleClient_BoxFileSystem() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	files, err := client.WithProject("sandbox").BoxFileSystem(context.Background(), "devbox")
	if err != nil {
		panic(err)
	}
	reader, err := files.Open(context.Background(), "/work/site/index.html", run9.OpenFileOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	_, _ = io.Copy(io.Discard, reader)
}

func ExampleClient_RunExecCapture() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	capture, err := project.RunExecCapture(context.Background(), "devbox", run9.ExecRequest{
		Command: []string{"/bin/sh", "-lc", "echo hello"},
	})
	if err != nil {
		panic(err)
	}

	_, _ = capture.Transcript, capture.Terminal
}

func ExampleClient_FollowBackgroundExec() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	follower := project.FollowBackgroundExec("exec-123")
	result, err := follower.Pump(context.Background(), 2*time.Second, run9.ExecOutputWriters{
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	if err != nil {
		panic(err)
	}

	_ = result.TerminalResult()
}

func ExampleClient_StartBackgroundExec() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		panic(err)
	}

	project := client.WithProject("sandbox")
	_, err = project.StartBackgroundExec(context.Background(), "devbox", run9.ExecRequest{
		Command:        []string{"/bin/sh", "-lc", "long task"},
		IdempotencyKey: "daemon-generation-42",
	})
	if err != nil {
		panic(err)
	}
}

func ExampleBackgroundExecPullOutput_WriteMergedOutput() {
	result := run9.BackgroundExecPullOutput{
		Events: []run9.BackgroundExecOutputEvent{
			{Type: run9.BackgroundExecOutputEventStdout, Data: []byte("hello\n")},
			{Type: run9.BackgroundExecOutputEventStderr, Data: []byte("warn\n")},
		},
	}

	if err := result.WriteMergedOutput(io.Discard, io.Discard); err != nil {
		panic(err)
	}
}

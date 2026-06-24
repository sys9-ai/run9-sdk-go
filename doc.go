/*
Package run9 provides the public Go SDK for the run9 control-plane API.

The package follows one consistent calling model:

  - Create one Client with NewClient by binding the API endpoint and API key once.
  - Use the base client for account, organization, shared-snap, and project-discovery APIs.
  - Use Client.WithProject to derive a project-scoped client for boxes, snaps, execs, archives, and secrets.
  - Pass request structs when a call has optional filters or mutable payload.
  - Use high-level helpers such as Client.RunExec and Client.FollowBackgroundExec when you want output plus terminal results.
  - Drop down to typed stream helpers such as ExecStream and ExecAttachSocket only when you need direct event control.

The module path is github.com/sys9-ai/run9-sdk-go, while the package name is run9:

	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		return
	}

	project := client.WithProject("default")
	_, _ = project.ListBoxes(context.Background(), run9.ListBoxesRequest{})

Non-2xx responses are returned as *Error. When the control plane returns a
structured error message, it is exposed through Error.Message.
*/
package run9

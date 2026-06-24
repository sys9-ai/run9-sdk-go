package run9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type errorPayload struct {
	Error string `json:"error"`
}

// Credentials are the org-scoped API key credentials used by run9 clients.
type Credentials struct {
	// AK is the access-key identifier.
	AK string
	// SK is the secret-key value paired with AK.
	SK string
}

// Error represents one run9 control-plane request failure.
type Error struct {
	// StatusCode is the HTTP status returned by the control plane.
	StatusCode int
	// Message is the structured control-plane error message when one is available.
	Message string
}

// Error returns the control-plane message when present.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("portal api request failed with status %d", e.StatusCode)
}

// Client is the public run9 control-plane HTTP client.
type Client struct {
	baseURL    string
	creds      Credentials
	http       *http.Client
	projectCID string
}

type requestOptions struct {
	query   map[string]string
	headers map[string]string
	body    any
	result  any
}

// NewClient creates one authenticated run9 control-plane HTTP client rooted at the given endpoint.
func NewClient(endpoint string, creds Credentials) (*Client, error) {
	baseURL, err := normalizeEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(creds.AK) == "" {
		return nil, errors.New("missing run9 access key")
	}
	if strings.TrimSpace(creds.SK) == "" {
		return nil, errors.New("missing run9 secret key")
	}
	return &Client{
		// Long-running control requests and streams must be bounded by caller contexts,
		// not by a shorter transport timeout inside the shared HTTP client.
		baseURL: baseURL,
		creds: Credentials{
			AK: strings.TrimSpace(creds.AK),
			SK: strings.TrimSpace(creds.SK),
		},
		http: &http.Client{},
	}, nil
}

// WithProject returns a shallow client clone pinned to one project CID.
func (c *Client) WithProject(projectCID string) *Client {
	if c == nil {
		return nil
	}
	clone := *c
	clone.projectCID = strings.TrimSpace(projectCID)
	return &clone
}

// WhoAmI loads the current authenticated user and organization identity.
func (c *Client) WhoAmI(ctx context.Context) (CurrentOrgIdentityView, error) {
	var view CurrentOrgIdentityView
	err := c.do(ctx, http.MethodGet, "/whoami", requestOptions{result: &view})
	return view, err
}

// CreateBox creates one project-scoped box.
func (c *Client) CreateBox(ctx context.Context, req CreateBoxRequest) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodPost, "/boxes", requestOptions{body: req, result: &view})
	return view, err
}

// ListBoxes lists project-scoped boxes with optional creator, label, and state filters.
func (c *Client) ListBoxes(ctx context.Context, req ListBoxesRequest) ([]BoxView, error) {
	query := map[string]string{}
	if strings.TrimSpace(req.Creator) != "" {
		query["creator"] = strings.TrimSpace(req.Creator)
	}
	if strings.TrimSpace(req.Label) != "" {
		query["label"] = strings.TrimSpace(req.Label)
	}
	if strings.TrimSpace(string(req.State)) != "" {
		query["state"] = strings.TrimSpace(string(req.State))
	}

	var views []BoxView
	err := c.doWorkspace(ctx, http.MethodGet, "/boxes", requestOptions{query: query, result: &views})
	return views, err
}

// GetBox loads one project-scoped box by ID.
func (c *Client) GetBox(ctx context.Context, boxID string) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodGet, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID)), requestOptions{result: &view})
	return view, err
}

// StopBox requests a graceful stop for one box.
func (c *Client) StopBox(ctx context.Context, boxID string) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodPost, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/stop", requestOptions{result: &view})
	return view, err
}

// DeleteBox deletes one box.
func (c *Client) DeleteBox(ctx context.Context, boxID string) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodDelete, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID)), requestOptions{result: &view})
	return view, err
}

// ImportSnap imports a snap from an image reference into the current project.
func (c *Client) ImportSnap(ctx context.Context, req ImportSnapRequest) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodPost, "/snaps/import", requestOptions{
		body: ImportSnapRequest{
			ImageRef: strings.TrimSpace(req.ImageRef),
		},
		result: &view,
	})
	return view, err
}

// ListSnaps lists project-scoped snaps with an optional attached filter.
func (c *Client) ListSnaps(ctx context.Context, req ListSnapsRequest) ([]SnapView, error) {
	query := map[string]string{}
	if req.Attached != nil {
		if *req.Attached {
			query["attached"] = "true"
		} else {
			query["attached"] = "false"
		}
	}

	var views []SnapView
	err := c.doWorkspace(ctx, http.MethodGet, "/snaps", requestOptions{query: query, result: &views})
	return views, err
}

// GetSnap loads one project-scoped snap by ID.
func (c *Client) GetSnap(ctx context.Context, snapID string) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodGet, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID)), requestOptions{result: &view})
	return view, err
}

// ForkSnap creates a writable child snap from an existing snap.
func (c *Client) ForkSnap(ctx context.Context, snapID string) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodPost, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID))+"/fork", requestOptions{result: &view})
	return view, err
}

// DeleteSnap deletes one snap.
func (c *Client) DeleteSnap(ctx context.Context, snapID string) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodDelete, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID)), requestOptions{result: &view})
	return view, err
}

// StartExecStream starts a streaming exec and returns one event reader.
func (c *Client) StartExecStream(ctx context.Context, boxID string, req ExecRequest) (*ExecStream, error) {
	cleanBoxID := strings.TrimSpace(boxID)
	resp, err := c.doWorkspaceRaw(ctx, http.MethodPost, "/boxes/"+url.PathEscape(cleanBoxID)+"/execs/stream", requestOptions{
		body: req,
		headers: map[string]string{
			"Accept":                  "application/x-ndjson",
			"X-Run9-Exec-Stream-Mode": "inline",
		},
	})
	if err != nil {
		return nil, err
	}
	return newExecStream(strings.TrimSpace(resp.Header.Get("X-Run9-Exec-ID")), resp.Body), nil
}

// RunExec starts one inline foreground exec, pumps its output, and returns the terminal result.
func (c *Client) RunExec(ctx context.Context, boxID string, req ExecRequest, writers ExecOutputWriters) (ExecTerminalResult, error) {
	stream, err := c.StartExecStream(ctx, boxID, req)
	if err != nil {
		return ExecTerminalResult{}, err
	}
	defer stream.Close()

	return stream.Pump(ctx, writers)
}

// StartExec starts one foreground exec and returns its initial view.
func (c *Client) StartExec(ctx context.Context, boxID string, req ExecRequest) (ExecView, error) {
	var view ExecView
	err := c.doWorkspace(ctx, http.MethodPost, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/execs", requestOptions{
		body:   req,
		result: &view,
	})
	return view, err
}

// UploadArchive uploads one tar archive into a box path.
func (c *Client) UploadArchive(ctx context.Context, boxID string, boxAbsPath string, source io.Reader) (RuntimeRequestView, error) {
	var view RuntimeRequestView
	err := c.doWorkspace(ctx, http.MethodPost, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/files/upload", requestOptions{
		query:   map[string]string{"box_abs_path": strings.TrimSpace(boxAbsPath)},
		headers: map[string]string{"Content-Type": "application/x-tar"},
		body:    source,
		result:  &view,
	})
	return view, err
}

// DownloadArchive downloads one box path as a tar archive.
func (c *Client) DownloadArchive(ctx context.Context, boxID string, boxAbsPath string) (io.ReadCloser, error) {
	resp, err := c.doWorkspaceRaw(ctx, http.MethodGet, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/files/download", requestOptions{
		query: map[string]string{
			"archive":      "tar",
			"box_abs_path": strings.TrimSpace(boxAbsPath),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) do(ctx context.Context, method string, path string, options requestOptions) error {
	resp, err := c.doRaw(ctx, method, path, options)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if options.result == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("portal api returned empty response body")
	}
	return json.Unmarshal(data, options.result)
}

func (c *Client) doRaw(ctx context.Context, method string, path string, options requestOptions) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, responseError(resp)
	}
	return resp, nil
}

func (c *Client) newRequest(ctx context.Context, method string, path string, options requestOptions) (*http.Request, error) {
	targetURL, err := requestURL(c.baseURL, path, options.query)
	if err != nil {
		return nil, err
	}
	body, err := requestBody(options.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.creds.AK, c.creds.SK)
	req.Header.Set("Accept", "application/json")
	for key, value := range options.headers {
		req.Header.Set(key, value)
	}
	if strings.TrimSpace(req.Header.Get("Content-Type")) == "" && options.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) doWorkspace(ctx context.Context, method string, path string, options requestOptions) error {
	workspacePath, err := c.workspacePath(path)
	if err != nil {
		return err
	}
	return c.do(ctx, method, workspacePath, options)
}

func (c *Client) doWorkspaceRaw(ctx context.Context, method string, path string, options requestOptions) (*http.Response, error) {
	workspacePath, err := c.workspacePath(path)
	if err != nil {
		return nil, err
	}
	return c.doRaw(ctx, method, workspacePath, options)
}

func requestURL(baseURL string, path string, query map[string]string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + path
	values := parsed.Query()
	for key, value := range query {
		values.Set(key, value)
	}
	parsed.RawQuery = values.Encode()
	parsed.Fragment = ""
	return parsed.String(), nil
}

func (c *Client) workspacePath(path string) (string, error) {
	cleanPath := "/" + strings.TrimLeft(strings.TrimSpace(path), "/")
	if strings.TrimSpace(c.projectCID) == "" {
		return "", errors.New("missing project cid: use client.WithProject(...) for project-scoped APIs")
	}
	return "/projects/" + url.PathEscape(c.projectCID) + "/workspace" + cleanPath, nil
}

func requestBody(body any) (io.Reader, error) {
	switch value := body.(type) {
	case nil:
		return nil, nil
	case io.Reader:
		return value, nil
	case []byte:
		return bytes.NewReader(value), nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}

func responseError(resp *http.Response) error {
	if resp == nil {
		return &Error{Message: "request failed"}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{StatusCode: resp.StatusCode, Message: err.Error()}
	}
	return responseBodyError(resp.StatusCode, resp.Status, data)
}

func responseBodyError(statusCode int, status string, body []byte) error {
	var payload errorPayload
	if len(body) > 0 && json.Unmarshal(body, &payload) == nil && strings.TrimSpace(payload.Error) != "" {
		return &Error{StatusCode: statusCode, Message: strings.TrimSpace(payload.Error)}
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = strings.TrimSpace(status)
	}
	return &Error{StatusCode: statusCode, Message: message}
}

func normalizeEndpoint(endpoint string) (string, error) {
	trimmed := strings.TrimSpace(endpoint)
	if trimmed == "" {
		return "", errors.New("missing run9 endpoint")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("parse endpoint: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid endpoint: %q", endpoint)
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("invalid endpoint: must not contain query or fragment: %q", endpoint)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed.String(), nil
}

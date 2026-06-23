package run9

import (
	"bytes"
	"context"
	"encoding/json"
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
	AK string
	SK string
}

// Error represents one run9 control-plane request failure.
type Error struct {
	StatusCode int
	Message    string
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
	http       *http.Client
	projectCID string
}

type requestOptions struct {
	query   map[string]string
	headers map[string]string
	body    any
	result  any
}

// NewClient creates one run9 control-plane HTTP client rooted at the given endpoint.
func NewClient(endpoint string) *Client {
	baseURL := strings.TrimRight(strings.TrimSpace(endpoint), "/")
	return &Client{
		// Long-running control requests and streams must be bounded by caller contexts,
		// not by a shorter transport timeout inside the shared HTTP client.
		baseURL: baseURL,
		http:    &http.Client{},
	}
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
func (c *Client) WhoAmI(ctx context.Context, creds Credentials) (CurrentOrgIdentityView, error) {
	var view CurrentOrgIdentityView
	err := c.do(ctx, http.MethodGet, "/whoami", creds, requestOptions{result: &view})
	return view, err
}

// CreateBox creates one project-scoped box.
func (c *Client) CreateBox(ctx context.Context, creds Credentials, req CreateBoxRequest) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/boxes"), creds, requestOptions{body: req, result: &view})
	return view, err
}

// Boxes lists project-scoped boxes with optional creator, label, and state filters.
func (c *Client) Boxes(ctx context.Context, creds Credentials, creator string, label string, state string) ([]BoxView, error) {
	query := map[string]string{}
	if strings.TrimSpace(creator) != "" {
		query["creator"] = strings.TrimSpace(creator)
	}
	if strings.TrimSpace(label) != "" {
		query["label"] = strings.TrimSpace(label)
	}
	if strings.TrimSpace(state) != "" {
		query["state"] = strings.TrimSpace(state)
	}

	var views []BoxView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/boxes"), creds, requestOptions{query: query, result: &views})
	return views, err
}

// Box loads one project-scoped box by ID.
func (c *Client) Box(ctx context.Context, creds Credentials, boxID string) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))), creds, requestOptions{result: &view})
	return view, err
}

// StopBox requests a graceful stop for one box.
func (c *Client) StopBox(ctx context.Context, creds Credentials, boxID string) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/stop"), creds, requestOptions{result: &view})
	return view, err
}

// RemoveBox deletes one box.
func (c *Client) RemoveBox(ctx context.Context, creds Credentials, boxID string) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodDelete, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))), creds, requestOptions{result: &view})
	return view, err
}

// ImportSnap imports a snap from an image reference into the current project.
func (c *Client) ImportSnap(ctx context.Context, creds Credentials, imageRef string) (SnapView, error) {
	var view SnapView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/snaps/import"), creds, requestOptions{
		body:   ImportSnapRequest{ImageRef: strings.TrimSpace(imageRef)},
		result: &view,
	})
	return view, err
}

// Snaps lists project-scoped snaps with an optional attached filter.
func (c *Client) Snaps(ctx context.Context, creds Credentials, attached string) ([]SnapView, error) {
	query := map[string]string{}
	if strings.TrimSpace(attached) != "" {
		query["attached"] = strings.TrimSpace(attached)
	}

	var views []SnapView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/snaps"), creds, requestOptions{query: query, result: &views})
	return views, err
}

// Snap loads one project-scoped snap by ID.
func (c *Client) Snap(ctx context.Context, creds Credentials, snapID string) (SnapView, error) {
	var view SnapView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/snaps/"+url.PathEscape(strings.TrimSpace(snapID))), creds, requestOptions{result: &view})
	return view, err
}

// ForkSnap creates a writable child snap from an existing snap.
func (c *Client) ForkSnap(ctx context.Context, creds Credentials, snapID string) (SnapView, error) {
	var view SnapView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/snaps/"+url.PathEscape(strings.TrimSpace(snapID))+"/fork"), creds, requestOptions{result: &view})
	return view, err
}

// RemoveSnap deletes one snap.
func (c *Client) RemoveSnap(ctx context.Context, creds Credentials, snapID string) (SnapView, error) {
	var view SnapView
	err := c.do(ctx, http.MethodDelete, c.workspacePath("/snaps/"+url.PathEscape(strings.TrimSpace(snapID))), creds, requestOptions{result: &view})
	return view, err
}

// ExecStream starts a streaming exec and returns the exec ID plus NDJSON body.
func (c *Client) ExecStream(ctx context.Context, creds Credentials, boxID string, req ExecBoxRequest) (string, io.ReadCloser, error) {
	cleanBoxID := strings.TrimSpace(boxID)
	resp, err := c.doRaw(ctx, http.MethodPost, c.workspacePath("/boxes/"+url.PathEscape(cleanBoxID)+"/execs/stream"), creds, requestOptions{
		body: req,
		headers: map[string]string{
			"Accept":                  "application/x-ndjson",
			"X-Run9-Exec-Stream-Mode": "inline",
		},
	})
	if err != nil {
		return "", nil, err
	}
	execID := strings.TrimSpace(resp.Header.Get("X-Run9-Exec-ID"))
	return execID, resp.Body, nil
}

// Exec starts one foreground exec and returns its initial view.
func (c *Client) Exec(ctx context.Context, creds Credentials, boxID string, req ExecBoxRequest) (ExecView, error) {
	var view ExecView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/execs"), creds, requestOptions{
		body:   req,
		result: &view,
	})
	return view, err
}

// UploadArchive uploads one tar archive into a box path.
func (c *Client) UploadArchive(ctx context.Context, creds Credentials, boxID string, boxAbsPath string, source io.Reader) (RuntimeRequestView, error) {
	var view RuntimeRequestView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/files/upload"), creds, requestOptions{
		query:   map[string]string{"box_abs_path": strings.TrimSpace(boxAbsPath)},
		headers: map[string]string{"Content-Type": "application/x-tar"},
		body:    source,
		result:  &view,
	})
	return view, err
}

// DownloadArchive downloads one box path as a tar archive.
func (c *Client) DownloadArchive(ctx context.Context, creds Credentials, boxID string, boxAbsPath string) (io.ReadCloser, error) {
	resp, err := c.doRaw(ctx, http.MethodGet, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))+"/files/download"), creds, requestOptions{
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

func (c *Client) do(ctx context.Context, method string, path string, creds Credentials, options requestOptions) error {
	resp, err := c.doRaw(ctx, method, path, creds, options)
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

func (c *Client) doRaw(ctx context.Context, method string, path string, creds Credentials, options requestOptions) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, creds, options)
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

func (c *Client) newRequest(ctx context.Context, method string, path string, creds Credentials, options requestOptions) (*http.Request, error) {
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
	req.SetBasicAuth(creds.AK, creds.SK)
	req.Header.Set("Accept", "application/json")
	for key, value := range options.headers {
		req.Header.Set(key, value)
	}
	if strings.TrimSpace(req.Header.Get("Content-Type")) == "" && options.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
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

func (c *Client) workspacePath(path string) string {
	cleanPath := "/" + strings.TrimLeft(strings.TrimSpace(path), "/")
	if strings.TrimSpace(c.projectCID) == "" {
		return cleanPath
	}
	return "/projects/" + url.PathEscape(c.projectCID) + "/workspace" + cleanPath
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

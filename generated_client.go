package run9

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	genclient "github.com/sys9-ai/run9-sdk-go/internal/generated/client"
)

var errEmptyResponseBody = errors.New("portal api returned empty response body")
var errNilGeneratedResponse = errors.New("generated client returned nil response")

func newGeneratedPortal(baseURL string, creds Credentials) (*genclient.Run9Portal, runtime.ClientAuthInfoWriter, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, nil, err
	}

	basePath := parsed.EscapedPath()
	if basePath == "" {
		basePath = "/"
	}

	httpClient := &http.Client{Transport: failureTransport{next: http.DefaultTransport}}
	transport := httptransport.NewWithClient(parsed.Host, basePath, []string{parsed.Scheme}, httpClient)
	auth := httptransport.BasicAuth(creds.AK, creds.SK)
	transport.Consumers[runtime.JSONMime] = nonEmptyResponseConsumer(transport.Consumers[runtime.JSONMime])
	transport.Consumers[runtime.TextMime] = nonEmptyResponseConsumer(transport.Consumers[runtime.TextMime])
	transport.Consumers[runtime.HTMLMime] = nonEmptyResponseConsumer(transport.Consumers[runtime.HTMLMime])
	return genclient.New(transport, nil), auth, nil
}

func nonEmptyResponseConsumer(inner runtime.Consumer) runtime.Consumer {
	return runtime.ConsumerFunc(func(reader io.Reader, data any) error {
		payload, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		if len(bytes.TrimSpace(payload)) == 0 {
			return errEmptyResponseBody
		}
		return inner.Consume(bytes.NewReader(payload), data)
	})
}

func remarshalJSON[To any](from any) (To, error) {
	var zero To
	if from == nil {
		return zero, nil
	}
	data, err := json.Marshal(from)
	if err != nil {
		return zero, err
	}
	var to To
	if err := json.Unmarshal(data, &to); err != nil {
		return zero, err
	}
	return to, nil
}

// failureTransport reports every failure status as *Error before the generated
// client reads the response. Generated operations only understand the statuses
// their route declares, and a proxy failure page is not JSON, so neither could
// otherwise reach callers with its status.
type failureTransport struct {
	next http.RoundTripper
}

// RoundTrip implements http.RoundTripper.
func (t failureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(req)
	if err != nil || resp.StatusCode < http.StatusBadRequest {
		return resp, err
	}
	return nil, responseError(resp)
}

// generatedError removes the transport wrapper from a failure status.
func generatedError(err error) error {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return err
}

func (c *Client) requireProjectCID() (string, error) {
	projectCID := strings.TrimSpace(c.projectCID)
	if projectCID == "" {
		return "", errors.New("missing project cid: use client.WithProject(...) for project-scoped APIs")
	}
	return projectCID, nil
}

func generatedResult[To any](result any, err error) (To, error) {
	var zero To
	if err != nil {
		return zero, generatedError(err)
	}

	payload, err := generatedPayload(result)
	if err != nil {
		return zero, err
	}
	return remarshalJSON[To](payload)
}

func generatedAction(_ any, err error) error {
	return generatedError(err)
}

func projectGeneratedResult[To any](c *Client, call func(projectCID string) (any, error)) (To, error) {
	var zero To
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return zero, err
	}
	return generatedResult[To](call(projectCID))
}

func projectGeneratedAction(c *Client, call func(projectCID string) (any, error)) error {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return err
	}
	return generatedAction(call(projectCID))
}

func generatedPayload(result any) (any, error) {
	if result == nil {
		return nil, errNilGeneratedResponse
	}

	value := reflect.ValueOf(result)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil, errNilGeneratedResponse
	}

	method := value.MethodByName("GetPayload")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return nil, fmt.Errorf("generated response %T has no GetPayload method", result)
	}
	return method.Call(nil)[0].Interface(), nil
}

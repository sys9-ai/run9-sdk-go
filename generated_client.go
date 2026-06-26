package run9

import (
	"bytes"
	"encoding/json"
	"errors"
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

func newGeneratedPortal(baseURL string, creds Credentials, httpClient *http.Client) (*genclient.Run9Portal, runtime.ClientAuthInfoWriter, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, nil, err
	}

	basePath := parsed.EscapedPath()
	if basePath == "" {
		basePath = "/"
	}

	transport := httptransport.NewWithClient(parsed.Host, basePath, []string{parsed.Scheme}, httpClient)
	auth := httptransport.BasicAuth(creds.AK, creds.SK)
	transport.DefaultAuthentication = auth
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

func generatedError(err error) error {
	if err == nil {
		return nil
	}

	coded, ok := err.(interface{ Code() int })
	if !ok {
		if errors.Is(err, errEmptyResponseBody) {
			return errEmptyResponseBody
		}
		return err
	}

	message := generatedErrorMessage(err)
	if message == "" {
		message = strings.TrimSpace(err.Error())
	}
	return &Error{
		StatusCode: coded.Code(),
		Message:    message,
	}
}

func generatedErrorMessage(err error) string {
	value := reflect.ValueOf(err)
	if !value.IsValid() {
		return ""
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ""
	}

	payload := value.FieldByName("Payload")
	if !payload.IsValid() {
		return ""
	}
	if payload.Kind() == reflect.Pointer {
		if payload.IsNil() {
			return ""
		}
		payload = payload.Elem()
	}
	if payload.Kind() != reflect.Struct {
		return ""
	}

	errorField := payload.FieldByName("Error")
	if !errorField.IsValid() || errorField.Kind() != reflect.String {
		return ""
	}
	return strings.TrimSpace(errorField.String())
}

func (c *Client) requireProjectCID() (string, error) {
	projectCID := strings.TrimSpace(c.projectCID)
	if projectCID == "" {
		return "", errors.New("missing project cid: use client.WithProject(...) for project-scoped APIs")
	}
	return projectCID, nil
}
